package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/kube-eip/internal/handler"
)

func LoadV1Api(app *gin.Engine) {
	v1 := app.Group("/v1")
	{
		v1.GET("/eipbindings", handler.ListEB)
	}
}
