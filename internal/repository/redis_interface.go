package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRepository Redis 仓储接口
type RedisRepository interface {
	// Set 设置缓存
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Get 获取缓存
	Get(ctx context.Context, key string, dest interface{}) error
	// Del 删除缓存
	Del(ctx context.Context, key string) error
	// Exists 检查是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// XAdd 向 Stream 追加消息
	XAdd(ctx context.Context, stream string, values map[string]interface{}) error
	// XReadGroup 消费者组读取消息
	XReadGroup(ctx context.Context, group, consumer string, streams []string, count int64, block time.Duration) ([]redis.XStream, error)
	// XAck 确认消息已消费
	XAck(ctx context.Context, stream, group string, ids ...string) error
	// XGroupCreate 创建消费者组
	XGroupCreate(ctx context.Context, stream, group, start string) error
	// XTrimMaxLen 截断 Stream 到指定长度
	XTrimMaxLen(ctx context.Context, stream string, maxLen int64) error
}
