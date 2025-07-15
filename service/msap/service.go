package msap

import (
	"StealthIMProxy/service/middleware"

	"github.com/gin-gonic/gin"
)

// SetService 定义路由
func SetService(router *gin.RouterGroup) {

	routerMsap := router.Group("/message")
	routerMsap.POST("/:groupid", middleware.AuthSession(), CheckGroup(), send)
	routerMsap.PATCH("/:groupid", middleware.AuthSession(), CheckGroup(), recall)
	routerMsap.GET("/:groupid", middleware.AuthSession(), CheckGroup(), sync)

}
