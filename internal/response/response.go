package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, JsonRsp{Code: 0, Message: "ok", Data: data})
}

func SuccessWithoutContent(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, JsonRsp{Code: 0, Message: "ok"})
}

func Failed(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, JsonRsp{Code: 1, Message: message})
}

func Forbidden(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, JsonRsp{Code: 1, Message: message})
}
