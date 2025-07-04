package persistence

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// 作品 repo
type WorksRepo struct {
	// DbCtx 是数据库上下文
	DbCtx *gorm.DB
}

// NewWorksRepo 构造一个新的 WorksRepo 实例
func NewWorksRepo(db *gorm.DB) *WorksRepo {
	return &WorksRepo{
		DbCtx: db,
	}
}

// GetTable 获取 works 表的上下文引用
func (r *WorksRepo) GetTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodels.Work{})
}
