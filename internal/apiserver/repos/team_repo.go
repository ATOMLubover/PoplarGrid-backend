package repos

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// TeamRepo 接口定义了团队仓库的基本操作
type TeamRepo interface {
	// SelectBasicPage 获取团队列表，支持分页
	SelectBasicPage(offset, limit int) ([]*dbmodels.Team, error)
}

// teamRepoImpl 是 TeamRepo 的实现
type teamRepoImpl struct {
	// handle 是数据库连接实例
	handle *gorm.DB
}

// NewTeamRepo 创建一个新的 TeamRepo 实例
func NewTeamRepo(db *gorm.DB) TeamRepo {
	return &teamRepoImpl{
		handle: db,
	}
}

// Table 限定当前操作的表名
func (r *teamRepoImpl) Table() *gorm.DB {
	return r.handle.Table(dbmodels.Team{}.TableName())
}

// SelectBasicPage 实现 TeamRepo 接口的 SelectBasicPage 方法
func (r *teamRepoImpl) SelectBasicPage(offset, limit int) ([]*dbmodels.Team, error) {
	var teams []*dbmodels.Team

	if err := r.Table().
		Select("id, name"). // 只选择需要的字段以提高性能
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&teams).
		Error; err != nil {
		return nil, err
	}

	return teams, nil
}
