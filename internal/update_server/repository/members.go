package repository

import (
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
)

// 成员 repo
type MembersRepo struct {
	// DbCtx 是数据库上下文
	DbCtx *gorm.DB
}

// NewMembersRepo 构造一个新的 MembersRepo 实例
func NewMembersRepo(db *gorm.DB) *MembersRepo {
	return &MembersRepo{
		DbCtx: db,
	}
}

// GetTable 获取 members 表的上下文引用
func (r *MembersRepo) GetTable() *gorm.DB {
	return r.DbCtx.Model(&models.Member{})
}

// BulkUpsert 批量更新或者插入成员
func (r *MembersRepo) BulkUpsert(inputMembers []*models.Member) error {
	if len(inputMembers) == 0 {
		return nil
	}

	// 为了实现高效的批量插入或更新，我们使用 GORM 的 Clauses(clause.OnConflict{})
	// 其基于 PostgreSQL 的 UPSERT 功能 UPSERT 支持，大幅减少了数据库的交互次数
	if err := r.GetTable().
		// Clauses(clause.OnConflict{
		// 	Columns:   []clause.Column{{Name: "moetran_id"}}, // 指定冲突字段为 moetran_id
		// 	DoNothing: true,
		// }).
		Create(&inputMembers).
		Error; err != nil {
		return err
	}

	return nil
}
