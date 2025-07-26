package models

import "gorm.io/gorm"

// Team 定义了团队的基本信息
type Team struct {
	BaseModel

	// 基本信息
	Name        string `gorm:"unique;size:128;not null"`
	Description string `gorm:"type:text"`

	// 龙译相关信息
	MoetranId string `gorm:"index;type:text"`
}

// 外部调用用的指针，避免每一次都创建
var emptyTeam = &Team{}

// GetTeam 获取一个空的 Team 实例
func GetTeam() *Team {
	return emptyTeam
}

// TableName 返回 Team 的表名
func (Team) TableName() string {
	return "teams"
}

// TeamSpec 定义了团队的查询条件
type TeamSpec struct {
	Id        *PKey
	Name      *string
	MoetranId *string
}

// Apply 将 TeamSpec 应用到 WHERE 子句上
func (s *TeamSpec) Apply(query *gorm.DB) {
	if s.Id != nil {
		query = query.Where("id = ?", *s.Id)
	}

	if s.Name != nil {
		query = query.Where("name = ?", *s.Name)
	}

	if s.MoetranId != nil {
		query = query.Where("moetran_id = ?", *s.MoetranId)
	}
}

// Insert 创建一个新的 Team 实例
func (*Team) Insert(hdl *gorm.DB, team *Team) error {
	if team == nil {
		return &InvalidParameterError{}
	}

	if team.Name == "" {
		return &LackOfRequiredFieldError{}
	}

	return hdl.
		Model(&Team{}).
		Create(team).Error
}

// Update 更新 Team 实例
func (*Team) Update(hdl *gorm.DB, team *Team) error {
	if team == nil {
		return &InvalidParameterError{}
	}

	if team.BaseModel.Id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Team{}).
		Updates(team).Error
}

// SelectFirst 查询一个 Team 实例
func (*Team) SelectFirst(
	hdl *gorm.DB, cnd *TeamSpec,
) (*Team, error) {
	var team Team

	query := hdl.Model(&Team{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if err := query.First(&team).Error; err != nil {
		return nil, err
	}

	return &team, nil
}

// SelectMany 查询多个 Team 实例
func (*Team) SelectMany(
	hdl *gorm.DB, cnd *TeamSpec,
	offset *int, limit *int,
) ([]*Team, error) {
	var teams []*Team

	query := hdl.Model(&Team{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Find(&teams).Error; err != nil {
		return nil, err
	}

	return teams, nil
}

// Delete 删除一个 Team 实例
func (*Team) Delete(hdl *gorm.DB, id PKey) error {
	if id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Team{}).
		Where("id = ?", id).
		Delete(&Team{}).Error
}
