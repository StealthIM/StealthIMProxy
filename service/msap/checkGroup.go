package msap

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"StealthIMProxy/conns"
	"StealthIMProxy/errorcode"

	pb "StealthIMProxy/StealthIM.GroupUser"

	"github.com/gin-gonic/gin"
)

func callGroupUser() (*pb.StealthIMGroupUserClient, error) {
	conn, err := conns.PoolGroup.ChooseConn()
	if err != nil {
		return nil, err
	}
	cli := pb.NewStealthIMGroupUserClient(conn)
	return &cli, nil
}

// CheckGroup 检查用户是否在群组中
func CheckGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var groupidStr string = c.Param("groupid")
		var uid int32 = int32(c.GetInt("uid"))
		groupid, err := strconv.ParseInt(groupidStr, 10, 64)

		if err != nil || groupid <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.ProxyBadRequest,
					"msg":  err.Error(),
				},
			})
			c.Abort()
			return
		}

		cli, err := callGroupUser()
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
		// 调用 GetGroupInfo 方法检查群组信息
		ret, err := (*cli).GetGroupInfo(ctx, &pb.GetGroupInfoRequest{
			Uid:     uid,
			GroupId: int32(groupid),
		})
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
		if ret.Result.Code != errorcode.Success {
			c.JSON(http.StatusOK, gin.H{
				"result": gin.H{
					"code": errorcode.GroupUserPermissionDenied,
					"msg":  ret.Result.Msg,
				},
			})
			c.Abort()
			return
		}

		c.Set("groupid", groupid)
		c.Next()
	}
}
