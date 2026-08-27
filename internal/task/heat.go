package task

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"blog/internal/model/entity"
	"blog/pkg/logger"
)

const (
	// HeatInterval 热度重算间隔
	HeatInterval = 1 * time.Hour
	// heatThreshold 热度变化死区阈值：|新热度 - 当前热度| 超过该值才写库，避免每小时全表更新
	heatThreshold = 3
	// heatExpr 热度计算公式：
	//   热度 = (基础热度10 + 浏览量*0.05 + 评论*6 + 点赞*1) * 时间系数
	//   时间系数 = 1 / (1 + 发布天数*0.02)，发布越久权重越低，永不归零
	//   基础热度10放在时间系数内：新内容发布即有非0热度可进「最热」，随发布天数衰减自然回落。
	heatExpr = "ROUND((10 + view_count * 0.05 + comment_count * 6 + like_count) * (1 / (1 + DATEDIFF(NOW(), published_at) * 0.02)))"
)

// StartHeatJob 启动热度定时重算任务：启动时立即算一次，之后按间隔循环。
// 返回 ticker，app 优雅关闭时调用 Stop 即可。
func StartHeatJob(db *gorm.DB) *time.Ticker {
	recomputeHeat(db)
	ticker := time.NewTicker(HeatInterval)
	go func() {
		for range ticker.C {
			recomputeHeat(db)
		}
	}()
	logger.Info("热度重算任务已启动", zap.Duration("interval", HeatInterval))
	return ticker
}

// recomputeHeat 按公式刷新 articles 表的 heat 列，只写「热度变化超过阈值」的行。
// 只重算已发布内容（status = 1）。
// 草稿/下架内容不参与计算，重新上架后会在下一次重算时恢复热度。
// 死区更新：WHERE 里用 ABS(新热度 - 当前热度) > heatThreshold 过滤，变化小的行直接跳过不写库。
func recomputeHeat(db *gorm.DB) {
	whereChanged := "ABS(" + heatExpr + " - hot) > ?"
	result := db.Model(&entity.Article{}).
		Where("status = ? AND published_at IS NOT NULL", 1).
		Where(whereChanged, heatThreshold).
		Update("hot", gorm.Expr(heatExpr))
	if result.Error != nil {
		logger.Error("热度重算失败(article)", zap.Error(result.Error))
		return
	}
	logger.Info("热度重算完成", zap.Int64("updated", result.RowsAffected))
}
