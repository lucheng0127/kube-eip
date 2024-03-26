package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/kube-eip/internal/router"
)

func Serve(port int) error {
	app := gin.Default()
	router.LoadV1Api(app)

	return app.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
