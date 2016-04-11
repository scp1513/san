package v1

import (
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	loginRouter(v1)
	regRouter(v1)
	serverRouter(v1)
}
