package repos

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// WorksetRepo 接口定义了工作集仓库的基本操作
type WorksetRepo interface {
	// SelectBasicPageIdDesc 按 ID 倒序获取工作集列表，支持分页
	SelectBasicPageIdDesc(teamId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.Workset, error)
}

// worksetRepo 是 WorksetRepo 的实现
type worksetRepo struct {
	handle *gorm.DB
}

// NewWorksetRepo 创建一个新的 WorksetRepo 实例
func NewWorksetRepo(db *gorm.DB) WorksetRepo {
	return &worksetRepo{
		handle: db,
	}
}

// Table 限定当前操作的表名
func (r *worksetRepo) Table() *gorm.DB {
	return r.handle.Table(dbmodels.Workset{}.TableName())
}

// SelectBasicPageIdDesc 实现 WorksetRepo 接口的 SelectBasicPageIdDesc 方法
func (r *worksetRepo) SelectBasicPageIdDesc(teamId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.Workset, error) {
	var worksets []*dbmodels.Workset

	if err := r.Table().
		Select("id, name, team_id"). // 只选择需要的字段以提高性能
		Where("team_id = ?", teamId).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&worksets).
		Error; err != nil {
		return nil, err
	}

	return worksets, nil
}
