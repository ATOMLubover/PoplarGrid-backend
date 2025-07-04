package persistence

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// 汉化组 repo，暂时应该用不上同步
type TeamsRepo struct {
	// DbCtx 是数据库上下文
	DbCtx *gorm.DB
}

// NewTeamsRepo 构造一个新的 TeamsRepo 实例
func NewTeamsRepo(db *gorm.DB) *TeamsRepo {
	return &TeamsRepo{
		DbCtx: db,
	}
}

// GetTable 获取 teams 表的上下文引用
func (r *TeamsRepo) GetTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodels.Team{})
}

// SelectNameAndId 获取 teams 表的名称和 ID
func (r *TeamsRepo) SelectNameAndId() ([]*dbmodels.Team, error) {
	var teams []*dbmodels.Team
	if err := r.GetTable().
		Select("id", "name", "moetran_id").
		Find(&teams).
		Error; err != nil {
		return nil, err
	}

	return teams, nil
}
