package api

import (
	"ddns-server/internal/config"
	"ddns-server/internal/dao"
	"ddns-server/internal/logic"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Record = aRecord{
	validator: validator.New(),
}

type aRecord struct {
	validator *validator.Validate
}

func (a *aRecord) List(ctx *gin.Context) {
	list, err := dao.Record.List()
	if err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Data: list,
	})
}

func (a *aRecord) Get(ctx *gin.Context) {
	domain := ctx.Param("domain")

	r, err := dao.Record.FindByDomain(domain)
	if err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Data: r,
	})
}

type RecordSetRequest struct {
	Domain          string `json:"domain" binding:"required"`
	Ipv4            string `json:"ipv4"`
	Ipv6            string `json:"ipv6"`
	ForceUpdateIpv4 bool   `json:"force_update_ipv4"`
	ForceUpdateIpv6 bool   `json:"force_update_ipv6"`
}

func (a *aRecord) Set(ctx *gin.Context) {
	request := new(RecordSetRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	if request.Ipv4 != "" {
		if err := a.validator.Var(request.Ipv4, "ipv4"); err != nil {
			ctx.JSON(http.StatusOK, ApiResponse{
				Errno:  1,
				Errmsg: err.Error(),
			})
			return
		}
	}

	if request.Ipv6 != "" {
		if err := a.validator.Var(request.Ipv6, "ipv6"); err != nil {
			ctx.JSON(http.StatusOK, ApiResponse{
				Errno:  1,
				Errmsg: err.Error(),
			})
			return
		}
	}

	err := dao.Record.Set(&config.Record{
		Domain: request.Domain,
		Ipv4:   request.Ipv4,
		Ipv6:   request.Ipv6,
	}, request.ForceUpdateIpv4, request.ForceUpdateIpv6)

	if err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	if err := logic.Dns.ReloadRecord(); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{})
}

func (a *aRecord) Delete(ctx *gin.Context) {
	domain := ctx.Param("domain")

	if err := dao.Record.Delete(domain); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	if err := logic.Dns.ReloadRecord(); err != nil {
		ctx.JSON(http.StatusOK, ApiResponse{
			Errno:  1,
			Errmsg: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{})
}
