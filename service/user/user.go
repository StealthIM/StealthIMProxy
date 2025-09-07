package user

import (
	"context"
	"net/http"
	"time"

	pb "StealthIMProxy/StealthIM.User"
	"StealthIMProxy/conns"
	"StealthIMProxy/errorcode"

	// 新增导入
	"github.com/gin-gonic/gin"
)

func call() (*pb.StealthIMUserClient, error) {
	conn, err := conns.PoolUser.ChooseConn()
	if err != nil {
		return nil, err
	}
	cli := pb.NewStealthIMUserClient(conn)
	return &cli, nil
}

func register(c *gin.Context) {
	var obj registerRequest
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).Register(ctx, &pb.RegisterRequest{
		Username:    obj.Username,
		Password:    obj.Password,
		Nickname:    obj.Nickname,
		Email:       obj.Email,
		PhoneNumber: obj.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
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

func login(c *gin.Context) {
	var obj loginRequest
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).Login(ctx, &pb.LoginRequest{
		Username: obj.Username,
		Password: obj.Password,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
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
		"session": ret.Session,
		"user_info": gin.H{
			"username":     ret.UserInfo.Username,
			"nickname":     ret.UserInfo.Nickname,
			"email":        ret.UserInfo.Email,
			"phone_number": ret.UserInfo.PhoneNumber,
			"vip":          ret.UserInfo.Vip,
			"create_time":  ret.UserInfo.CreateTime,
		},
	})
}

func delete(c *gin.Context) {
	cli, err := call()
	var uid int32 = int32(c.GetInt("uid"))
	if uid <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ProxyAuthFailed,
				"msg":  "uid is invalid",
			},
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).Logout(ctx, &pb.LogoutRequest{
		UserId: uid,
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

func getUserInfo(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).GetUserInfo(ctx, &pb.GetUserInfoRequest{
		UserId: uid,
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
		"user_info": gin.H{
			"username":     ret.UserInfo.Username,
			"nickname":     ret.UserInfo.Nickname,
			"email":        ret.UserInfo.Email,
			"phone_number": ret.UserInfo.PhoneNumber,
			"vip":          ret.UserInfo.Vip,
			"create_time":  ret.UserInfo.CreateTime,
		},
	})
}

func getUserPublicInfo(c *gin.Context) {
	var username string = c.Param("username")
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).GetOtherUserInfo(ctx, &pb.GetOtherUserInfoRequest{
		Username: username,
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
		"user_info": gin.H{
			"nickname": ret.UserInfo.Nickname,
			"vip":      ret.UserInfo.Vip,
		},
	})
}

func changePassword(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	var obj changePasswordRequest
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).ChangePassword(ctx, &pb.ChangePasswordRequest{
		UserId:      uid,
		NewPassword: obj.Password,
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

func changeNickname(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	var obj changeNicknameRequest
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).ChangeNickname(ctx, &pb.ChangeNicknameRequest{
		UserId:      uid,
		NewNickname: obj.Nickname,
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

func changeEmail(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	var obj changeEmailRequest
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).ChangeEmail(ctx, &pb.ChangeEmailRequest{
		UserId:   uid,
		NewEmail: obj.Email,
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

func changePhoneNumber(c *gin.Context) {
	var uid int32 = int32(c.GetInt("uid"))
	var obj changePhoneNumberRequest
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ret, err := (*cli).ChangePhoneNumber(ctx, &pb.ChangePhoneNumberRequest{
		UserId:         uid,
		NewPhoneNumber: obj.PhoneNumber,
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
