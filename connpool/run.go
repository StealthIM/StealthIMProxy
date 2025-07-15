package connpool

import (
	"errors"
	"math/rand"

	"google.golang.org/grpc"
)

// ChooseConn 随机选择链接
func (p *Pool) ChooseConn() (*grpc.ClientConn, error) {
	if len(*p.conns) == 0 {
		return nil, errors.New("No available connections")
	}
	for range 5 {
		conntmp := (*p.conns)[rand.Intn(len(*p.conns))]
		if conntmp != nil {
			return conntmp, nil
		}
	}
	return nil, errors.New("No available connections")
}
