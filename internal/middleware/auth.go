package middleware

import (
	"strings"

	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/jwt"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	// ContextUserID 用户 ID 上下文键
	ContextUserID = "user_id"
	// ContextUsername 用户名 上下文键
	ContextUsername = "username"
)

// Auth JWT 认证中间件
func Auth(authRepo repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请提供认证令牌")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "令牌格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 检查 token 是否在黑名单中
		if authRepo != nil {
			exists, _ := authRepo.IsTokenBlacklisted(token)
			if exists {
				response.Unauthorized(c, "token 已失效")
				c.Abort()
				return
			}
		}

		// 解析 Token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Error(c, errors.CodeTokenExpired, err.Error())
			} else {
				response.Error(c, errors.CodeInvalidToken, err.Error())
			}
			c.Abort()
			return
		}

		// 检查用户状态
		if authRepo != nil {
			user, err := authRepo.FindByID(claims.GetUserID())
			if err != nil || user == nil {
				response.Unauthorized(c, "用户不存在")
				c.Abort()
				return
			}
			if user.Status != 1 {
				response.Error(c, errors.CodeUserDisabled, "账号已被禁用")
				c.Abort()
				return
			}
		}

		// 将用户信息存入上下文
		c.Set(ContextUserID, claims.GetUserID())
		c.Set(ContextUsername, claims.GetUsername())

		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get(ContextUserID); exists {
		return userID.(uint)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(ContextUsername); exists {
		return username.(string)
	}
	return ""
}

// AdminAuth 后台管理 JWT 认证中间件（查 author 表）
func AdminAuth(authRepo repository.AuthRepository, authorRepo repository.AuthorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请提供认证令牌")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "令牌格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 检查 token 是否在黑名单中
		if authRepo != nil {
			exists, _ := authRepo.IsTokenBlacklisted(token)
			if exists {
				response.Unauthorized(c, "token 已失效")
				c.Abort()
				return
			}
		}

		// 解析 Token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Error(c, errors.CodeTokenExpired, err.Error())
			} else {
				response.Error(c, errors.CodeInvalidToken, err.Error())
			}
			c.Abort()
			return
		}

		// 检查作者状态（查 author 表）
		author, err := authorRepo.FindByID(claims.GetUserID())
		if err != nil || author == nil {
			response.Unauthorized(c, "作者不存在")
			c.Abort()
			return
		}

		// 将作者信息存入上下文
		c.Set(ContextUserID, author.ID)
		c.Set(ContextUsername, author.Account)

		c.Next()
	}
}

// OptionalAuth 可选的 JWT 认证中间件
func OptionalAuth(authRepo repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]

		// 检查 token 是否在黑名单中
		if authRepo != nil {
			exists, _ := authRepo.IsTokenBlacklisted(token)
			if exists {
				c.Next()
				return
			}
		}

		claims, err := jwt.ParseToken(token)
		if err == nil {
			c.Set(ContextUserID, claims.GetUserID())
			c.Set(ContextUsername, claims.GetUsername())
		}

		c.Next()
	}
}
