package services

import (
	"log/slog"
	"poplargrid/internal/api_server/crawler"
)

// CrawlerService 定义了爬虫服务的接口
type CrawlerService interface {
	// AutoUpdateAll 自动递归地更新当前用户所有汉化组的项目信息
	AutoUpdateAll(userID uint, moetranAuth string) Err
}

// crawlerServiceImpl 实现了 CrawlerService 接口
type crawlerServiceImpl struct {
	crawler crawler.Crawler
	logger  *slog.Logger
}

// NewCrawlerService 创建一个新的 CrawlerService 实例
func NewCrawlerService(
	crw crawler.Crawler,
	lgr *slog.Logger,
) CrawlerService {
	return &crawlerServiceImpl{
		crawler: crw,
		logger:  lgr,
	}
}

// AutoUpdateAll 自动递归地更新当前用户所有汉化组的项目信息
func (s *crawlerServiceImpl) AutoUpdateAll(userID uint, moetranAuth string) Err {
	// 调用爬虫的 AutoUpdateAll 方法
	if err := s.crawler.RecurseUpdate(userID, moetranAuth); err != nil {
		s.logger.Error("AutoUpdateAll 调用爬虫服务失败",
			slog.Any("error", err))
		return newSrvError(ErrMoetranAPIFailure, "自动拉取尨译失败: "+err.Error())
	}

	return nil
}
