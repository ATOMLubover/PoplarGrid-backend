package repository

import (
	"poplargrid/internal/shared/models"

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
	return r.DbCtx.Model(&models.Team{})
}

// Select 获取所有汉化组的详细信息
func (r *TeamsRepo) Select() ([]*models.Team, error) {
	var teams []*models.Team
	if err := r.GetTable().
		Find(&teams).
		Error; err != nil {
		return nil, err
	}

	return teams, nil
}
