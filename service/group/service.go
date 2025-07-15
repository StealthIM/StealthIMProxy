package group

import (
	"StealthIMProxy/service/middleware"

	"github.com/gin-gonic/gin"
)

// SetService 定义路由
func SetService(router *gin.RouterGroup) {
	routerGrp := router.Group("/group")
	routerGrp.GET("/", middleware.AuthSession(), getByUID)
	routerGrp.GET("/:groupid", middleware.AuthSession(), getInfo)
	routerGrp.GET("/:groupid/public", middleware.AuthSession(), getPublicInfo)
	routerGrp.POST("/", middleware.AuthSession(), create)
	routerGrp.POST("/:groupid/join", middleware.AuthSession(), join)
	routerGrp.POST("/:groupid/invite", middleware.AuthSession(), invite)
	routerGrp.PUT("/:groupid/:username", middleware.AuthSession(), setUserType)
	routerGrp.DELETE("/:groupid/:username", middleware.AuthSession(), kickUser)
	routerGrp.PATCH("/:groupid/name", middleware.AuthSession(), changeName)
	routerGrp.PATCH("/:groupid/password", middleware.AuthSession(), changePassword)
}
