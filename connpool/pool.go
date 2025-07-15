package connpool

import (
	"StealthIMProxy/config"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Pool 连接池
type Pool struct {
	conns          *[]*grpc.ClientConn
	mainlock       sync.RWMutex
	checkAliveFunc func(*grpc.ClientConn) (bool, error)
	name           string
	cfg            *config.NodeConfig
}

func createConn(connID int, host string, port int, conns *[]*grpc.ClientConn, name string) {
	log.Printf("[%s]Connect %d", name, connID+1)
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))
	if conn == nil {
		log.Printf("[%s]Connect %d Error %v\n", name, connID+1, err)
		(*conns)[connID] = nil
		return
	}
	if err != nil {
		log.Printf("[%s]Connect %d Error %v\n", name, connID+1, err)
		(*conns)[connID] = nil
		return
	}
	(*conns)[connID] = conn
}

func checkAlive(connID int, host string, port int, checkAlivefn func(*grpc.ClientConn) (bool, error), conns *[]*grpc.ClientConn, name string, mainlock *sync.RWMutex) {
	if len(*conns) <= connID {
		return
	}
	for {
		if len(*conns) <= connID {
			return
		}
		mainlock.RLock()
		if (*conns)[connID] != nil {
			online, err := checkAlivefn((*conns)[connID])
			if err == nil && online {
				mainlock.RUnlock()
				continue
			}
		}
		createConn(connID, host, port, conns, name)
		mainlock.RUnlock()
		time.Sleep(5 * time.Second)
	}
}

func (p *Pool) initConns() {
	defer func() {
		p.mainlock.Lock()
		for _, conn := range *p.conns {
			conn.Close()
		}
		p.mainlock.Unlock()
	}()
	log.Printf("[%s]Init Conns\n", p.name)
	for {
		time.Sleep(time.Second * 1)
		var lenTmp = len((*p.conns))
		if lenTmp < p.cfg.ConnNum {
			log.Printf("[%s]Create Conn %d\n", p.name, lenTmp+1)
			p.mainlock.Lock()
			*p.conns = append(*p.conns, nil)
			p.mainlock.Unlock()
			go checkAlive(lenTmp, p.cfg.Host, p.cfg.Port, p.checkAliveFunc, p.conns, p.name, &p.mainlock)
		} else if lenTmp > p.cfg.ConnNum {
			log.Printf("[%s]Delete Conn %d\n", p.name, lenTmp)
			p.mainlock.Lock()
			(*p.conns)[lenTmp-1].Close()
			*p.conns = (*p.conns)[:lenTmp-1]
			p.mainlock.Unlock()
		} else {
			time.Sleep(time.Second * 5)
		}
	}
}

// NewPool 创建连接池
func NewPool(name string, cfg config.NodeConfig, checkAliveFunc func(*grpc.ClientConn) (bool, error)) *Pool {
	pool := &Pool{
		conns:          &[]*grpc.ClientConn{},
		mainlock:       sync.RWMutex{},
		checkAliveFunc: checkAliveFunc,
		name:           name,
		cfg:            &cfg,
	}
	go pool.initConns()
	return pool
}
