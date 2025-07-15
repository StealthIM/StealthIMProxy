package group

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	pb "StealthIMProxy/StealthIM.GroupUser"
	"StealthIMProxy/conns"
	"StealthIMProxy/errorcode"

	// 新增导入
	"github.com/gin-gonic/gin"
)

func call() (*pb.StealthIMGroupUserClient, error) {
	conn, err := conns.PoolGroup.ChooseConn()
	if err != nil {
		return nil, err
	}
	cli := pb.NewStealthIMGroupUserClient(conn)
	return &cli, nil
}

func getByUID(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).GetGroupsByUID(ctx, &pb.GetGroupsByUIDRequest{
		Uid: uid,
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": ret.Result.Code,
				"msg":  ret.Result.Msg,
			},
		})
		return
	}
	var grp = make([]int32, len(ret.Groups))
	copy(grp, ret.Groups)
	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"code": errorcode.Success,
			"msg":  "",
		},
		"groups": grp,
	})
}

func getPublicInfo(c *gin.Context) {
	var groupidStr string = c.Param("groupid")
	groupid, err := strconv.ParseInt(groupidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).GetGroupPublicInfo(ctx, &pb.GetGroupPublicInfoRequest{
		GroupId: int32(groupid),
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
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
		"name":      ret.Name,
		"create_at": fmt.Sprintf("%d", ret.CreatedAt),
	})
}

func getInfo(c *gin.Context) {
	var uid = int32(c.GetInt("uid"))
	var groupidStr string = c.Param("groupid")
	groupid, err := strconv.ParseInt(groupidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).GetGroupInfo(ctx, &pb.GetGroupInfoRequest{
		Uid:     uid,
		GroupId: int32(groupid),
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": ret.Result.Code,
				"msg":  ret.Result.Msg,
			},
		})
		return
	}
	var members = make([]gin.H, 0)
	for _, v := range ret.Members {
		members = append(members, gin.H{
			"name": v.Name,
			"type": int(v.Type),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"code": errorcode.Success,
			"msg":  "",
		},
		"members": ret.Members,
	})
}

func create(c *gin.Context) {
	var obj createGroupRequest
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}
	var uid int32 = int32(c.GetInt("uid"))
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).CreateGroup(ctx, &pb.CreateGroupRequest{
		Uid:  uid,
		Name: obj.Name,
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
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
		"groupid": ret.GroupId,
	})
}

func join(c *gin.Context) {
	var uid = int32(c.GetInt("uid"))
	var groupidStr string = c.Param("groupid")
	groupid, err := strconv.ParseInt(groupidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
	}
	var obj joinRequest
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).JoinGroup(ctx, &pb.JoinGroupRequest{
		Password: obj.Password,
		Uid:      uid,
		GroupId:  int32(groupid),
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
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

func invite(c *gin.Context) {
	var uid = int32(c.GetInt("uid"))
	var groupidStr string = c.Param("groupid")
	groupid, err := strconv.ParseInt(groupidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
	}
	var obj inviteRequest
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).InviteGroup(ctx, &pb.InviteGroupRequest{
		GroupId:  int32(groupid),
		Uid:      uid,
		Username: obj.Username,
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
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

func setUserType(c *gin.Context) {
	var uid = int32(c.GetInt("uid"))
	var groupidStr string = c.Param("groupid")
	var username string = c.Param("username")
	groupid, err := strconv.ParseInt(groupidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
	}
	var obj setUserTypeRequest
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}
	if obj.Type != int32(pb.MemberType_owner) && obj.Type != int32(pb.MemberType_member) && obj.Type != int32(pb.MemberType_manager) {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  "Type is invalid",
			},
		})
		return
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).SetUserType(ctx, &pb.SetUserTypeRequest{
		GroupId:  int32(groupid),
		Uid:      uid,
		Username: username,
		Type:     pb.MemberType(obj.Type),
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
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

func kickUser(c *gin.Context) {
	var uid = int32(c.GetInt("uid"))
	var username string = c.Param("username")
	var groupidStr string = c.Param("groupid")
	groupid, err := strconv.ParseInt(groupidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).KickUser(ctx, &pb.KickUserRequest{
		GroupId:  int32(groupid),
		Uid:      uid,
		Username: username,
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
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

func changePassword(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	var groupidStr string = c.Param("groupid")
	groupid, err := strconv.ParseInt(groupidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
	}
	var obj changePasswordRequest
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).ChangeGroupPassword(ctx, &pb.ChangeGroupPasswordRequest{
		Password: obj.Password,
		GroupId:  int32(groupid),
		Uid:      uid,
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
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

func changeName(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	var groupidStr string = c.Param("groupid")
	groupid, err := strconv.ParseInt(groupidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
	}
	var obj changeNameRequest
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  err.Error(),
			},
		})
		return
	}
	cli, err := call()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).ChangeGroupName(ctx, &pb.ChangeGroupNameRequest{
		GroupId: int32(groupid),
		Name:    obj.Name,
		Uid:     uid,
	})
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusBadRequest, gin.H{
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
