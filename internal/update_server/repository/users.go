package repository

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

type UsersRepo struct {
	dbCtx *gorm.DB
}

// NewUsersRepo 创建一个新的 UsersRepo 实例
func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{
		dbCtx: db,
	}
}

// GetTable 获取 users 表的上下文引用
func (r *UsersRepo) GetTable() *gorm.DB {
	return r.dbCtx.Model(&dbmodels.User{})
}

// BulkUpsert 批量更新或插入用户
func (r *UsersRepo) BulkUpsert(users []*dbmodels.User) error {
	if len(users) == 0 {
		return nil
	}

	// 使用 GORM 的 Clauses(clause.OnConflict{}) 实现批量插入或更新
	if err := r.GetTable().
		// Clauses(clause.OnConflict{
		// 	Columns:   []clause.Column{{Name: "moetran_id"}}, // 指定冲突字段为 moetran_id
		// 	DoUpdates: clause.Assignments(map[string]any{"nickname": gorm.Expr("EXCLUDED.nickname")}),
		// }).
		Create(&users).
		Error; err != nil {
		return err
	}

	return nil
}
