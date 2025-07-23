package msap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	pb "StealthIMProxy/StealthIM.MSAP"
	"StealthIMProxy/conns"
	"StealthIMProxy/errorcode"

	"github.com/gin-gonic/gin"
)

func call() (*pb.StealthIMMSAPClient, error) {
	conn, err := conns.PoolMSAPW.ChooseConn()
	if err != nil {
		return nil, err
	}
	cli := pb.NewStealthIMMSAPClient(conn)
	return &cli, nil
}
func callR() (*pb.StealthIMMSAPSyncClient, error) {
	conn, err := conns.PoolMSAPR.ChooseConn()
	if err != nil {
		return nil, err
	}
	cli := pb.NewStealthIMMSAPSyncClient(conn)
	return &cli, nil
}

func send(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	var obj sendRequest
	groupid := c.GetInt64("groupid")
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}

	cli, err := call()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).SendMessage(ctx, &pb.SendMessageRequest{
		Uid:     uid,
		Msg:     obj.Msg,
		Groupid: int64(groupid),
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": ret.Result.Code,
				"msg":  ret.Result.Msg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"code": errorcode.Success,
			"msg":  "",
		},
	})
}

func recall(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	var obj recallRequest

	groupid := c.GetInt64("groupid")
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}
	msgid, err := strconv.ParseInt(obj.MsgID, 10, 64)

	if err != nil || msgid <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).RecallMessage(ctx, &pb.RecallMessageRequest{
		Uid:     uid,
		Msgid:   msgid,
		Groupid: int64(groupid),
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": ret.Result.Code,
				"msg":  ret.Result.Msg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"code": errorcode.Success,
			"msg":  "",
		},
	})
}

func sync(c *gin.Context) {
	// 设置SSE相关的HTTP头
	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 从请求中获取groupid和last_msgid
	groupID := c.GetInt64("groupid")
	lastMsgIDStr := c.Query("msgid")

	if groupID == 0 {
		jsonData, _ := json.Marshal(gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  "GroupID is empty",
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		return
	}
	msgid, err := strconv.ParseInt(lastMsgIDStr, 10, 64)

	if err != nil || msgid < 0 {
		jsonData, _ := json.Marshal(gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  "bad request",
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		return
	}

	cli, err := callR()

	if err != nil {
		jsonData, _ := json.Marshal(gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Failed to connect to gRPC server",
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		return
	}

	// 创建gRPC流请求
	req := &pb.SyncMessageRequest{
		Groupid:   groupID,
		LastMsgid: msgid,
	}

	// 调用gRPC服务端流方法
	stream, err := (*cli).SyncMessage(c.Request.Context(), req)
	if err != nil {
		jsonData, _ := json.Marshal(gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Failed to connect to gRPC server",
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		return
	}

	go func() {
		for {
			time.Sleep(30 * time.Second)
			select {
			case <-c.Request.Context().Done():
				return
			default:
				jsonData, err := json.Marshal(gin.H{
					"result": gin.H{
						"code": errorcode.Success,
						"msg":  "",
					},
					"msg": []gin.H{},
				})
				fmt.Fprintf(w, "data: %s\n\n", jsonData)
				flusher.Flush()
				if err != nil {
					return
				}
			}
		}
	}()

	// 循环从gRPC流中接收消息并发送到SSE客户端
	for {
		select {
		case <-c.Request.Context().Done():
			return
		default:
			msg, err := stream.Recv()
			if err != nil {
				return
			}

			var jsonData []byte
			// 将gRPC响应转换为JSON
			if msg.Result.Code != errorcode.Success {
				jsonData, err = json.Marshal(gin.H{
					"result": gin.H{
						"code": msg.Result.Code,
						"msg":  msg.Result.Msg,
					},
				})
				fmt.Fprintf(w, "data: %s\n\n", jsonData)
				flusher.Flush() // 立即发送数据到客户端
				return
			}
			msgs := make([]gin.H, len(msg.Msg))
			for i, m := range msg.Msg {
				msgs[i] = gin.H{
					"msgid":   fmt.Sprintf("%d", m.Msgid), // 处理int64类型
					"uid":     m.Uid,
					"groupid": fmt.Sprintf("%d", m.Groupid), // 处理int64类型
					"msg":     m.Msg,
					"type":    int(m.Type),
					"hash":    m.Hash,
					"time":    fmt.Sprintf("%d", m.Time), // 处理int64类型
				}
			}
			jsonData, err = json.Marshal(gin.H{
				"result": gin.H{
					"code": errorcode.Success,
					"msg":  "",
				},
				"msg": msgs,
			})
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		}
	}
}
