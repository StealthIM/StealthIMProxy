package file

import (
	"context"
	"net/http"
	"time"

	"StealthIMProxy/conns"
	"StealthIMProxy/errorcode"

	pb "StealthIMProxy/StealthIM.Session"

	"github.com/gin-gonic/gin"
)

func callSession() (*pb.StealthIMSessionClient, error) {
	conn, err := conns.PoolSession.ChooseConn()
	if err != nil {
		return nil, err
	}
	cli := pb.NewStealthIMSessionClient(conn)
	return &cli, nil
}

// AuthSessionForWS 从query验证session
func AuthSessionForWS() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Query("authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.ProxyAuthFailed,
					"msg":  "Authorization query is required",
				},
			})
			c.Abort() // 终止请求，不再执行后续的中间件和路由
			return
		}

		tokenString := authHeader

		cli, err := callSession()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.ServerInternalNetworkError,
					"msg":  err.Error(),
				},
			})
			c.Abort()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ret, err := (*cli).Get(ctx, &pb.GetRequest{
			Session: tokenString,
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.ProxyAuthFailed,
					"msg":  err.Error(),
				},
			})
			c.Abort()
			return
		}
		if ret.Result.Code != errorcode.Success {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.ProxyAuthFailed,
					"msg":  ret.Result.Msg,
				},
			})
			c.Abort()
			return
		}
		if ret.Uid <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.ProxyAuthFailed,
					"msg":  "Session not found",
				},
			})
			c.Abort()
			return
		}
		var uid int = int(ret.Uid)

		c.Set("uid", uid)

		c.Next()
	}
}
