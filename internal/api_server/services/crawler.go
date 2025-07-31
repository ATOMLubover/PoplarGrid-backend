package services

import "poplargrid/internal/api_server/crawler"

// CrawlerService 定义了爬虫服务的接口
type CrawlerService interface {
	// AutoUpdateAll 自动递归地更新当前用户所有汉化组的项目信息
	AutoUpdateAll(moetranAuth string) error
}

// crawlerServiceImpl 实现了 CrawlerService 接口
type crawlerServiceImpl struct {
	crawler crawler.Crawler
}

// NewCrawlerService 创建一个新的 CrawlerService 实例
func NewCrawlerService(crw crawler.Crawler) CrawlerService {
	return &crawlerServiceImpl{
		crawler: crw,
	}
}

// AutoUpdateAll 自动递归地更新当前用户所有汉化组的项目信息
func (s *crawlerServiceImpl) AutoUpdateAll(moetranAuth string) error {
	// 调用爬虫的 AutoUpdateAll 方法
	return s.crawler.AutoUpdateAll(moetranAuth)
}
