package repos

import "gorm.io/gorm"

// Repo 接口定义了项目仓库的基本操作
type Repo interface {
	// 获取数据库上下文句柄
	GetHandle() *gorm.DB
}
