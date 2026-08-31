package dashboard

import (
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 概览控制器
type Controller struct {
	svc service.DashboardService
}

// NewController 创建概览控制器
func NewController(svc service.DashboardService) *Controller {
	return &Controller{svc: svc}
}

// adminGetDashboard GET /api/v1/admin/dashboard
func (ctrl *Controller) adminGetDashboard(c *gin.Context) {
	result, err := ctrl.svc.GetDashboard()
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}
