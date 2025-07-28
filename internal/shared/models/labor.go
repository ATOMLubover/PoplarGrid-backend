package models

import "gorm.io/gorm"

// Labor 定义了 project 中一个 member 的基本信息
type Labor struct {
	BaseModel

	// 分工，按掩码存储，因为没有查询需求
	LaborMask uint32 `gorm:"not null;default:0"`

	// 所属项目外键
	ProjectId PKey     `gorm:"not null;index"`
	FkProject *Project `gorm:"foreignKey:ProjectId"`

	// 所属成员外键
	MemberId PKey    `gorm:"not null;index"`
	FkMember *Member `gorm:"foreignKey:MemberId"`
}

// 外部调用用的指针，避免每一次都创建
var emptyLabor = &Labor{}

// GetLabor 获取一个空的 Labor 实例
func GetLabor() *Labor {
	return emptyLabor
}

// TableName 返回 Labor 的表名
func (*Labor) TableName() string {
	return "labors"
}

// LaborSpec 定义了 Labor 的查询条件
type LaborSpec struct {
	Id *PKey

	// 所属项目外键
	ProjectId *PKey

	// 所属成员外键
	MemberId *PKey
}

// Apply 将 LaborSpec 应用为 WHERE 子句
func (s *LaborSpec) Apply(query *gorm.DB) {
	if s.Id != nil {
		query = query.Where("id = ?", *s.Id)
	}

	if s.ProjectId != nil {
		query = query.Where("project_id = ?", *s.ProjectId)
	}

	if s.MemberId != nil {
		query = query.Where("member_id = ?", *s.MemberId)
	}
}

// LaborFields 定义了 Labor 的预加载字段
type LaborFields struct {
	Id   bool
	Time bool

	ProjectId     bool
	ProjectFields *ProjectFields // 需要 Preload 时指定

	MemberId     bool
	MemberFields *MemberFields // 需要 Preload 时指定
}

// Apply 将 LaborFields 应用到 SELECT 子句
func (l *LaborFields) Apply(query *gorm.DB) {
	fields := make([]string, 0)

	if l.Id {
		fields = append(fields, "id")
	}
	if l.Time {
		fields = append(fields, "created_at", "updated_at", "deleted_at")
	}

	if l.ProjectId {
		fields = append(fields, "project_id")
		if l.ProjectFields != nil {
			query = query.Preload("FkProject", func(db *gorm.DB) *gorm.DB {
				l.ProjectFields.Apply(db)
				return db
			})
		}
	}
	if l.MemberId {
		fields = append(fields, "member_id")
		if l.MemberFields != nil {
			query = query.Preload("FkMember", func(db *gorm.DB) *gorm.DB {
				l.MemberFields.Apply(db)
				return db
			})
		}
	}

	query.Select(fields)
}

// Insert 插入 Labor 实例到数据库
func (*Labor) Insert(hdl *gorm.DB, labor *Labor) error {
	if labor == nil {
		return &InvalidParameterError{}
	}

	return hdl.
		Model(&Labor{}).
		Create(labor).Error
}

// SelectFirst 查询一个 Labor 实例
func (*Labor) SelectFirst(
	hdl *gorm.DB, cnd *LaborSpec, fields *LaborFields,
) (*Labor, error) {
	var labor Labor

	query := hdl.Model(&Labor{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if err := query.
		First(&labor).Error; err != nil {
		return nil, err
	}

	return &labor, nil
}

// SelectMany 查询多个 Labor 实例
func (*Labor) SelectMany(
	hdl *gorm.DB, cnd *LaborSpec, fields *LaborFields,
) ([]*Labor, error) {
	var labors []*Labor

	query := hdl.Model(&Labor{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if err := query.
		Find(&labors).Error; err != nil {
		return nil, err
	}

	return labors, nil
}

// Update 更新 Labor 实例
func (*Labor) Update(hdl *gorm.DB, labor *Labor) error {
	if labor == nil {
		return &InvalidParameterError{}
	}

	if labor.BaseModel.Id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Labor{}).
		Where("id = ?", labor.BaseModel.Id).
		Updates(labor).Error
}

// Delete 删除 Labor 实例
func (*Labor) Delete(hdl *gorm.DB, laborId PKey) error {
	if laborId == 0 {
		return &InvalidParameterError{}
	}

	return hdl.
		Model(&Labor{}).
		Where("id = ?", laborId).
		Delete(&Labor{}).Error
}
