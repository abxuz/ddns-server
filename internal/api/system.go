package api

import (
	"ddns-server/internal/config"
	"ddns-server/internal/dao"
	"ddns-server/internal/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

var System = aSystem{}

type aSystem struct {
}

type SystemGetResponse struct {
	Web *config.Web `json:"web,omitempty"`
	Dns *config.Dns `json:"dns,omitempty"`
}

func (a *aSystem) Get(ctx *gin.Context) {
	response := new(SystemGetResponse)

	var err error
	response.Web, err = dao.Web.Get()
	if err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	response.Dns, err = dao.Dns.Get()
	if err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Data: response,
	})
}

type SystemSetRequest struct {
	Web *config.Web `json:"web"`
	Dns *config.Dns `json:"dns"`
}

func (a *aSystem) Set(ctx *gin.Context) {
	request := new(SystemSetRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	if err := dao.Dns.Set(request.Dns); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	if err := logic.Dns.ReloadService(); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	if err := dao.Web.Set(request.Web); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	if err := logic.Web.ReloadService(); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{})
}
