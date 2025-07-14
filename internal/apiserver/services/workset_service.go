package services

import (
	"log/slog"
	"poplargrid/internal/apiserver/repos"
	"poplargrid/internal/shared/dbmodels"
)

// WorksetService 接口定义了作品集服务的基本操作
type WorksetService interface {
	// GetBasicPageIdDesc 获取作品集列表，按 ID 倒序
	GetBasicPage(pageSerial, pageSize int) ([]*dbmodels.Workset, error)
}

// worksetServiceImpl 是 WorksetService 的实现
type worksetServiceImpl struct {
	worksetRepo repos.WorksetRepo
	logger      *slog.Logger
}
