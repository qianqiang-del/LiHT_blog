package service

import dto "blog/internal/model/dto/response"

// DashboardService 概览服务接口
type DashboardService interface {
	// GetDashboard 获取概览数据
	GetDashboard() (*dto.DashboardDTO, error)
}
