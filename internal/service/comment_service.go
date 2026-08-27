package service

import (
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"

	"gorm.io/gorm"
)

// commentService 评论服务实现
type commentService struct {
	repo repository.CommentRepository
	db   *gorm.DB
}

// NewCommentService 创建评论服务
func NewCommentService(repo repository.CommentRepository, db *gorm.DB) CommentService {
	return &commentService{repo: repo, db: db}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(articleID, userID uint, nickname, content string, parentID *uint) (*dto.CommentItem, error) {
	comment := &entity.Comment{
		ArticleID: articleID,
		ParentID:  parentID,
		UserID:    &userID,
		Nickname:  nickname,
		Content:   content,
		Status:    1,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, errors.New(errors.CodeInternalError, "创建评论失败")
	}

	return &dto.CommentItem{
		ID:        comment.ID,
		Nickname:  comment.Nickname,
		Content:   comment.Content,
		LikeCount: 0,
		Liked:     false,
		CreatedAt: comment.CreatedAt,
	}, nil
}

// ListComments 获取文章的一级评论列表
func (s *commentService) ListComments(articleID uint, page, size int, userID *uint) (*dto.CommentListResponse, error) {
	offset := (page - 1) * size

	comments, total, err := s.repo.ListByArticleID(articleID, offset, size)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询评论列表失败")
	}

	// 统计回复数
	commentIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		commentIDs = append(commentIDs, c.ID)
	}
	replyCountMap, err := s.repo.CountReplies(commentIDs)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询回复数失败")
	}

	// 检查当前用户是否点赞
	likedMap := make(map[uint]bool)
	if userID != nil {
		for _, c := range comments {
			liked, _ := s.repo.HasLiked(s.db, c.ID, *userID)
			likedMap[c.ID] = liked
		}
	}

	// 组装 DTO
	list := make([]dto.CommentItem, 0, len(comments))
	for _, c := range comments {
		item := dto.CommentItem{
			ID:         c.ID,
			Nickname:   c.Nickname,
			Content:    c.Content,
			LikeCount:  c.LikeCount,
			Liked:      likedMap[c.ID],
			CreatedAt:  c.CreatedAt,
			ReplyCount: replyCountMap[c.ID],
		}
		list = append(list, item)
	}

	return &dto.CommentListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// ListReplies 获取评论的回复列表
func (s *commentService) ListReplies(commentID uint, page, size int, userID *uint) (*dto.CommentListResponse, error) {
	offset := (page - 1) * size

	replies, total, err := s.repo.ListByParentID(commentID, offset, size)
	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "查询回复列表失败")
	}

	// 检查当前用户是否点赞
	likedMap := make(map[uint]bool)
	if userID != nil {
		for _, r := range replies {
			liked, _ := s.repo.HasLiked(s.db, r.ID, *userID)
			likedMap[r.ID] = liked
		}
	}

	// 组装 DTO
	list := make([]dto.CommentItem, 0, len(replies))
	for _, r := range replies {
		item := dto.CommentItem{
			ID:        r.ID,
			Nickname:  r.Nickname,
			Content:   r.Content,
			LikeCount: r.LikeCount,
			Liked:     likedMap[r.ID],
			CreatedAt: r.CreatedAt,
		}
		list = append(list, item)
	}

	return &dto.CommentListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// LikeComment 点赞/取消点赞评论
func (s *commentService) LikeComment(commentID, userID uint) (*dto.LikeResponse, error) {
	var liked bool
	var likeCount int

	err := s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.repo.HasLiked(tx, commentID, userID)
		if err != nil {
			return err
		}

		if exists {
			// 取消点赞（DeleteLike 内部已处理 like_count -1）
			if err := s.repo.DeleteLike(tx, commentID, userID); err != nil {
				return err
			}
			liked = false
		} else {
			// 点赞（CreateLike 内部已处理 like_count +1）
			if err := s.repo.CreateLike(tx, commentID, userID); err != nil {
				return err
			}
			liked = true
		}

		count, err := s.repo.GetLikeCount(tx, commentID)
		if err != nil {
			return err
		}
		likeCount = count
		return nil
	})

	if err != nil {
		return nil, errors.New(errors.CodeInternalError, "操作失败")
	}

	return &dto.LikeResponse{Liked: liked, LikeCount: likeCount}, nil
}
