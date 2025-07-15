package file

import (
	"StealthIMProxy/service/middleware"

	"github.com/gin-gonic/gin"
)

// SetService 定义路由
func SetService(router *gin.RouterGroup) {
	routerGrp := router.Group("/file")
	routerGrp.GET("/", AuthSessionForWS(), upload)
	routerGrp.POST("/:filehash", middleware.AuthSession(), getInfo)
	routerGrp.GET("/:filehash", middleware.AuthSession(), download)
}
