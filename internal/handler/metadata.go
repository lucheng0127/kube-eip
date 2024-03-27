package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucheng0127/kube-eip/internal/response"
	"github.com/lucheng0127/kube-eip/pkg/utils/metadata"
)

type ListFilter struct {
	EIp    string `form:"exip"`
	IIp    string `form:"inip"`
	Status string `form:"status"`
}

func ListEB(c *gin.Context) {
	filter := new(ListFilter)
	err := c.ShouldBindQuery(filter)
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	mds, err := metadata.ListMD()
	if err != nil {
		response.Failed(c, err.Error())
		return
	}

	var data []*metadata.EipMetadata
	if filter == nil {
		data = mds
	} else {
		for _, md := range mds {
			if filter.EIp != "" {
				if md.ExternalIP != filter.EIp {
					continue
				}
			}

			if filter.IIp != "" {
				if md.InternalIP != filter.IIp {
					continue
				}
			}

			if filter.Status == "succeed" {
				if md.Status != metadata.MD_STATUS_FINISHED {
					continue
				}
			} else if filter.Status == "failed" {
				if md.Status != metadata.MD_STATUS_FAILED {
					continue
				}
			}

			data = append(data, md)
		}
	}

	response.Success(c, data)
	return
}
