package handlers

import (
	"poplargrid/internal/gateway_server/dtos"

	"github.com/kataras/iris/v12"
)

// 统筹向尨译转发的 handler
type TransferHandler struct {
}

// GetHealthCheck godoc
// @Summary 检查 transfer 路由的连通性
// @Tags transfer
// @Produce json
// @Success 200 {object} map[string]string
// @Router /try_ping [get]
func (h *TransferHandler) HealthCheck(ctx iris.Context) {
	ctx.JSON(dtos.JsonKV{
		"status": "OK",
		"server": "Gateway Server",
	})
}

func (h *TransferHandler) TransferToMoetran(ctx iris.Context) {

}
