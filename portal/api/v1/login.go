package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/scp1513/san/portal/models/login"
)

// @SubApi 登陆 [/login]
func loginRouter(r *gin.RouterGroup) {
	router := r.Group("/login")
	router.POST("/account", loginAccount)
	router.POST("/verify", loginVerify)
}

func loginAccount(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, login.Login(c.Request.Form))
}

func loginVerify(c *gin.Context) {
	c.Request.ParseForm()
	c.JSON(200, login.LoginVerify(c.Request.Form))
}
