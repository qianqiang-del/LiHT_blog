package upload

import (
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 上传控制器
type Controller struct {
	svc service.UploadService
}

// NewController 创建上传控制器
func NewController(svc service.UploadService) *Controller {
	return &Controller{svc: svc}
}

// upload POST /api/v1/upload
func (ctrl *Controller) upload(c *gin.Context) {
	if ctrl.svc == nil {
		response.InternalError(c, "上传服务未配置")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 限制文件大小 10MB
	if header.Size > 10*1024*1024 {
		response.BadRequest(c, "文件大小不能超过 10MB")
		return
	}

	url, err := ctrl.svc.Upload(c.Request.Context(), file, header.Filename)
	if err != nil {
		response.InternalError(c, "上传文件失败")
		return
	}

	response.Success(c, gin.H{"url": url})
}
