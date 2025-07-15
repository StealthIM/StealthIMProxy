package service

import (
	"StealthIMProxy/config"
	"StealthIMProxy/service/file"
	"StealthIMProxy/service/group"
	"StealthIMProxy/service/msap"
	"StealthIMProxy/service/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

// Start 启动服务
func Start(cfg config.Config) {
	engine := gin.Default()

	router := engine.Group("/api/v1")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, StealthIM!",
		})
	})
	user.SetService(router)
	group.SetService(router)
	msap.SetService(router)
	file.SetService(router)

	engine.Run(fmt.Sprintf("%s:%d", cfg.Proxy.Host, cfg.Proxy.Port))
}
