package user

import (
	"StealthIMProxy/service/middleware"

	"github.com/gin-gonic/gin"
)

// SetService 定义路由
func SetService(router *gin.RouterGroup) {
	routerUsr := router.Group("/user")
	routerUsr.POST("/register", register)
	routerUsr.POST("/", login)
	routerUsr.DELETE("/", middleware.AuthSession(), delete)
	routerUsr.GET("/", middleware.AuthSession(), getUserInfo)
	routerUsr.PATCH("/password", middleware.AuthSession(), changePassword)
	routerUsr.PATCH("/nickname", middleware.AuthSession(), changeNickname)
	routerUsr.PATCH("/email", middleware.AuthSession(), changeEmail)
	routerUsr.PATCH("/phone", middleware.AuthSession(), changePhoneNumber)
	routerUsr.GET("/:username", middleware.AuthSession(), getUserPublicInfo)
}
