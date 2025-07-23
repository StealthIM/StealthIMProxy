package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"StealthIMProxy/conns"
	"StealthIMProxy/errorcode"

	pb "StealthIMProxy/StealthIM.Session"

	"github.com/gin-gonic/gin"
)

func call() (*pb.StealthIMSessionClient, error) {
	conn, err := conns.PoolSession.ChooseConn()
	if err != nil {
		return nil, err
	}
	cli := pb.NewStealthIMSessionClient(conn)
	return &cli, nil
}

// AuthSession 是一个 Gin 中间件
// 它从 Authorization 头中提取 Token，并将其解析为 UID，然后存储到 Context 中。
func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头中获取 Authorization 字段
		authHeader := c.GetHeader("Authorization")

		// 2. 检查 Authorization 头是否存在且格式正确 (例如 "Bearer <token>")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.ProxyAuthFailed,
					"msg":  "Authorization header is required",
				},
			})
			c.Abort() // 终止请求，不再执行后续的中间件和路由
			return
		}

		// 假设 Token 格式是 "Bearer <token_string>"
		// 我们需要分割字符串来获取实际的 token_string
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.ProxyAuthFailed,
					"msg":  "Authorization header format must be Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		cli, err := call()
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
