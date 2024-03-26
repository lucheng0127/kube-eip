package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/kube-eip/internal/response"
	"github.com/lucheng0127/kube-eip/pkg/utils/metadata"
)

type ListFilter struct {
	eip    string
	iip    string
	status string
}

func ListEB(c *gin.Context) {
	filter := new(ListFilter)
	err := c.ShouldBind(filter)
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	mds, err := metadata.ListMD()
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	// TODO(shawnlu): Add filter

	response.Success(c, mds)
	return
}
