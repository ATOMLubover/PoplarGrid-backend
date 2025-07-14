package repos

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// WorksetRepo 接口定义了工作集仓库的基本操作
type WorksetRepo interface {
	// SelectBasicPageIdDesc 按 ID 倒序获取工作集列表，支持分页
	SelectBasicPageIdDesc(offset, limit int) ([]*dbmodels.Workset, error)
}

// worksetRepo 是 WorksetRepo 的实现
type worksetRepo struct {
	db *gorm.DB // 假设 dbmodels.DB 是你的数据库连接类型
}

// NewWorksetRepo 创建一个新的 WorksetRepo 实例
func NewWorksetRepo(db *gorm.DB) WorksetRepo {
	return &worksetRepo{
		db: db,
	}
}

// Table 限定当前操作的表名
func (r *worksetRepo) Table() *gorm.DB {
	return r.db.Table(dbmodels.Workset{}.TableName())
}

// SelectBasicPageIdDesc 实现 WorksetRepo 接口的 SelectBasicPageIdDesc 方法
func (r *worksetRepo) SelectBasicPageIdDesc(offset, limit int) ([]*dbmodels.Workset, error) {
	var worksets []*dbmodels.Workset

	if err := r.Table().
		Select("id, name, team_id"). // 只选择需要的字段以提高性能
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&worksets).
		Error; err != nil {
		return nil, err
	}

	return worksets, nil
}
