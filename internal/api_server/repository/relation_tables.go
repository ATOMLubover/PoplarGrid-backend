package repository

import (
	"poplargrid/internal/shared/dbmodel"

	"gorm.io/gorm"
)

// 由于关联表不是复杂的实体表，故统一管理入口
type RelationTables struct {
	DbCtx *gorm.DB
}

// NewRelationTables 构造一个新的 RelationTables 实例
func NewRelationTables(db *gorm.DB) *RelationTables {
	return &RelationTables{
		DbCtx: db,
	}
}

// 获取成员偏好表的上下文引用
func (r *RelationTables) GetMemberPreferenceTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodel.MemberPreference{})
}

// 获取成员分工表的上下文引用
func (r *RelationTables) GetProjectLaborDivisionTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodel.ProjectLaborDivision{})
}

// 获取作品 tag 表的上下文引用
func (r *RelationTables) GetProjectTagTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodel.ProjectTag{})
}
