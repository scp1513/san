package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/scp1513/san/portal/models/server"
)

// @SubApi 登陆 [/server]
func serverRouter(r *gin.RouterGroup) {
	router := r.Group("/server")
	router.POST("/list", serverList)
	router.POST("/id", serverID)
	router.POST("/release", serverRelease)
	router.POST("/stress", serverStress)
	router.POST("/addr", serverAddr)
}

func serverList(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, server.List(c.Request.Form))
}

func serverID(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, server.GetID(c.Request.Form))
}

func serverRelease(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, server.Release(c.Request.Form))
}

func serverStress(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, server.Stress(c.Request.Form))
}

func serverAddr(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, server.GetAddr(c.Request.Form))
}
