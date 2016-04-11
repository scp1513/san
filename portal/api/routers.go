package api

import (
	"github.com/gin-gonic/gin"

	"github.com/scp1513/san/portal/api/v1"
)

// SetupRouters 注册路由表
func SetupRouters(router *gin.Engine) {
	root := router.Group("/portal")

	v1.Setup(root)

	// catch no router
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, "404! page not found!")
	})
}
