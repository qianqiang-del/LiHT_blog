package stream

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"blog/internal/repository"
	"blog/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Redis Stream 配置
const (
	StreamKey     = "stream:view_count"
	ConsumerGroup = "cg:view_count"
	ConsumerName  = "consumer:1"

	BatchMax      = 50
	FlushInterval = 5 * time.Second
	StreamMaxLen  = 10000
)

// Consumer 浏览量异步消费者
type Consumer struct {
	redisRepo repository.RedisRepository
	db        *gorm.DB

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// NewConsumer 创建消费者实例
func NewConsumer(redisRepo repository.RedisRepository, db *gorm.DB) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		redisRepo: redisRepo,
		db:        db,
		ctx:       ctx,
		cancel:    cancel,
		done:      make(chan struct{}),
	}
}

// Start 创建消费者组并启动消费循环
func (c *Consumer) Start() {
	if c.redisRepo == nil {
		logger.Warn("Stream Consumer: Redis 未初始化，不启动")
		close(c.done)
		return
	}
	if err := c.ensureGroup(); err != nil {
		logger.Warn("Stream Consumer: 创建消费者组失败，消费循环中将自动重试", zap.Error(err))
	}
	go c.consume()
	logger.Info("Stream Consumer 已启动",
		zap.String("stream", StreamKey),
		zap.String("group", ConsumerGroup),
	)
}

// ensureGroup 确保 Consumer Group 存在（MKSTREAM: Stream 不存在时自动创建）
func (c *Consumer) ensureGroup() error {
	return c.redisRepo.XGroupCreate(c.ctx, StreamKey, ConsumerGroup, "0")
}

// Stop 优雅停止消费者
func (c *Consumer) Stop() {
	c.cancel()
	<-c.done
	logger.Info("Stream Consumer 已停止")
}

// consume 主循环
func (c *Consumer) consume() {
	defer close(c.done)
	defer c.restartOnPanic()

	counts := make(map[uint]int64)
	var msgIDs []string
	lastFlush := time.Now()

	for {
		select {
		case <-c.ctx.Done():
			c.flush(counts, msgIDs)
			return
		default:
		}

		c.pullAndMerge(counts, &msgIDs)

		n := len(msgIDs)
		if n >= BatchMax || (n > 0 && time.Since(lastFlush) >= FlushInterval) {
			if err := c.flush(counts, msgIDs); err != nil {
				logger.Warn("flush 失败，等待重试", zap.Error(err))
				time.Sleep(2 * time.Second)
				continue
			}
			counts = make(map[uint]int64)
			msgIDs = nil
			lastFlush = time.Now()
		}
	}
}

// pullAndMerge 拉取并合并消息
func (c *Consumer) pullAndMerge(counts map[uint]int64, msgIDs *[]string) {
	streams, err := c.redisRepo.XReadGroup(c.ctx,
		ConsumerGroup, ConsumerName,
		[]string{StreamKey, ">"}, int64(BatchMax), FlushInterval,
	)
	if err != nil {
		if errors.Is(err, redis.Nil) || c.ctx.Err() != nil {
			return
		}
		// NOGROUP: Consumer Group 不存在（Redis 重启后丢失），自动重建
		if isNoGroupErr(err) {
			logger.Warn("Stream Consumer: Consumer Group 不存在，正在重建...")
			if createErr := c.ensureGroup(); createErr != nil {
				logger.Warn("Stream Consumer: 重建 Consumer Group 失败", zap.Error(createErr))
			}
			return
		}
		logger.Warn("Stream Consumer: XReadGroup 失败", zap.Error(err))
		time.Sleep(time.Second)
		return
	}

	for _, s := range streams {
		for _, m := range s.Messages {
			id, ok := parseMsg(m)
			if !ok {
				_ = c.redisRepo.XAck(c.ctx, StreamKey, ConsumerGroup, m.ID)
				continue
			}
			counts[id]++
			*msgIDs = append(*msgIDs, m.ID)
		}
	}
}

// isNoGroupErr 判断错误是否为 NOGROUP（Consumer Group 不存在）
func isNoGroupErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "NOGROUP")
}

// flush 批量写入 MySQL
func (c *Consumer) flush(counts map[uint]int64, msgIDs []string) error {
	if len(msgIDs) == 0 {
		return nil
	}

	for id, n := range counts {
		if err := c.db.Exec("UPDATE articles SET view_count = view_count + ? WHERE id = ?", n, id).Error; err != nil {
			logger.Error("Stream Consumer: UPDATE 失败",
				zap.Uint("id", id),
				zap.Int64("count", n),
				zap.Error(err),
			)
			return err
		}
	}

	_ = c.redisRepo.XAck(c.ctx, StreamKey, ConsumerGroup, msgIDs...)
	_ = c.redisRepo.XTrimMaxLen(c.ctx, StreamKey, StreamMaxLen)
	return nil
}

// restartOnPanic panic 时自动重启
func (c *Consumer) restartOnPanic() {
	if r := recover(); r != nil {
		logger.Error("Stream Consumer panic，重启", zap.Any("panic", r))
		c.ctx, c.cancel = context.WithCancel(context.Background())
		c.done = make(chan struct{})
		go c.consume()
	}
}

// parseMsg 解析消息，提取文章 ID
func parseMsg(m redis.XMessage) (uint, bool) {
	idStr, _ := m.Values["id"].(string)
	n, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(n), true
}
