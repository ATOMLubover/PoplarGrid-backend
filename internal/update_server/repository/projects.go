package repository

import (
	"fmt"
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// 项目 repo，所有与本地数据库有关的项目均为不缩写的 Project
type ProjectsRepo struct {
	// DbCtx 是数据库上下文
	DbCtx *gorm.DB
}

// NewProjectsRepo 构造一个新的 WorksRepo 实例
func NewProjectsRepo(db *gorm.DB) *ProjectsRepo {
	return &ProjectsRepo{
		DbCtx: db,
	}
}

// GetTable 获取 works 表的上下文引用
func (r *ProjectsRepo) GetTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodels.Project{})
}

// BulkUpsert 批量更新或者插入项目
func (r *ProjectsRepo) BulkUpsert(inputProjects []*dbmodels.Project) error {
	if len(inputProjects) == 0 {
		return nil
	}

	// 为了实现高效的批量插入或更新，我们使用 GORM 的 Clauses(clause.OnConflict{})
	// 其基于 PostgreSQL 的 UPSERT 功能 UPSERT 支持，大幅减少了数据库的交互次数
	if err := r.GetTable().
		// Clauses(clause.OnConflict{
		// 	Columns: []clause.Column{{
		// 		Name: "moetran_id", // 指定冲突字段为 moetran_id
		// 	}},
		// 	DoUpdates: clause.Assignments(map[string]any{
		// 		"title":     gorm.Expr("EXCLUDED.title"),
		// 		"legacy_id": gorm.Expr("EXCLUDED.legacy_id"),
		// 	}),
		// }).
		Create(&inputProjects).
		Error; err != nil {
		return fmt.Errorf("批量 Upsert Projects 失败: %w", err)
	}

	return nil
}
