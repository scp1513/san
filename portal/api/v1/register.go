package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/scp1513/san/portal/models/register"
)

// @SubApi 登陆 [/reg]
func regRouter(r *gin.RouterGroup) {
	router := r.Group("/reg")
	router.POST("/vistor", regVistor)
	router.POST("/upgrade", regUpgrade)
	router.POST("/doreg", regDoreg)
}

func regVistor(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, register.Visitor(c.Request.Form))
}

func regUpgrade(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, register.Upgrade(c.Request.Form))
}

func regDoreg(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, register.Reg(c.Request.Form))
}
