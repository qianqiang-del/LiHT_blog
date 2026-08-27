package comment

import (
	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createComment POST /api/v1/articles/:id/comments
func (ctrl *Controller) createComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	nickname := middleware.GetUsername(c)

	result, err := ctrl.svc.CreateComment(uint(articleID), userID, nickname, req.Content, req.ParentID)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// Controller 评论控制器
type Controller struct {
	svc service.CommentService
}

// NewController 创建评论控制器
func NewController(svc service.CommentService) *Controller {
	return &Controller{svc: svc}
}

// listComments GET /api/v1/articles/:id/comments
func (ctrl *Controller) listComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	var req request.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}

	// OptionalAuth：未登录 userID 为 nil
	var userID *uint
	if uid := middleware.GetUserID(c); uid > 0 {
		userID = &uid
	}

	result, err := ctrl.svc.ListComments(uint(articleID), req.Page, req.Size, userID)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// likeComment POST /api/v1/comments/:id/like
func (ctrl *Controller) likeComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	result, err := ctrl.svc.LikeComment(uint(commentID), userID)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}

// listReplies GET /api/v1/comments/:id/replies
func (ctrl *Controller) listReplies(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}

	var req request.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}

	// OptionalAuth：未登录 userID 为 nil
	var userID *uint
	if uid := middleware.GetUserID(c); uid > 0 {
		userID = &uid
	}

	result, err := ctrl.svc.ListReplies(uint(commentID), req.Page, req.Size, userID)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, result)
}
