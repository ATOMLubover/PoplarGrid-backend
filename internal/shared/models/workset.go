package models

import "gorm.io/gorm"

// Workset 定义了工作集的基本信息
type Workset struct {
	BaseModel

	// 基本信息
	Name        string `gorm:"unique;size:128;not null"`
	Description string `gorm:"type:text"`

	// 用于项目的组内自增序列
	// 其 trigger 直接在数据库中实现
	ProjectSequenceName string `gorm:"uniqueIndex;type:text;not null"`

	// 龙译相关信息
	MoetranId string `gorm:"index;type:text"`

	// 所属团队外键
	TeamId PKey  `gorm:"not null;index"`
	FkTeam *Team `gorm:"foreignKey:TeamId"`
}

// 外部调用用的指针，避免每一次都创建
var emptyWorkset = &Workset{}

// GetWorkset 获取一个空的 Workset 实例
func GetWorkset() *Workset {
	return emptyWorkset
}

// TableName 返回 Workset 的表名
func (*Workset) TableName() string {
	return "worksets"
}

// WorksetSpec 定义了工作集的查询条件
type WorksetSpec struct {
	Id *PKey

	Name *string

	MoetranId *string

	TeamId *PKey
}

// Apply 将 WorksetSpec 应用到 WHERE 子句
func (s *WorksetSpec) Apply(query *gorm.DB) {
	if s.Id != nil {
		query = query.Where("id = ?", *s.Id)
	}

	if s.Name != nil {
		query = query.Where("name = ?", *s.Name)
	}

	if s.MoetranId != nil {
		query = query.Where("moetran_id = ?", *s.MoetranId)
	}

	if s.TeamId != nil {
		query = query.Where("team_id = ?", *s.TeamId)
	}
}

// Insert 创建一个新的 Workset 实例
func (*Workset) Insert(hdl *gorm.DB, workset *Workset) error {
	if workset == nil {
		return &InvalidParameterError{}
	}

	if workset.Name == "" {
		return &LackOfRequiredFieldError{}
	}

	return hdl.
		Model(&Workset{}).
		Create(workset).Error
}

// Update 更新 Workset 实例
func (*Workset) Update(hdl *gorm.DB, workset *Workset) error {
	if workset == nil {
		return &InvalidParameterError{}
	}

	if workset.BaseModel.Id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Workset{}).
		Updates(workset).Error
}

// SelectFirst 查询一个工作集
func (*Workset) SelectFirst(
	hdl *gorm.DB, cnd *WorksetSpec,
) (*Workset, error) {
	var workset Workset

	query := hdl.Model(&Workset{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if err := query.
		First(&workset).Error; err != nil {
		return nil, err
	}

	return &workset, nil
}

// SelectMany 查询多个工作集
func (*Workset) SelectMany(
	hdl *gorm.DB, cnd *WorksetSpec,
	offset *int, limit *int,
) ([]*Workset, error) {
	var worksets []*Workset

	query := hdl.Model(&Workset{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.
		Find(&worksets).Error; err != nil {
		return nil, err
	}

	return worksets, nil
}

// Delete 删除一个工作集
func (*Workset) Delete(hdl *gorm.DB, id PKey) error {
	if id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Workset{}).
		Where("id = ?", id).
		Delete(Workset{}).Error
}
