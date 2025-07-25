package models

import "gorm.io/gorm"

// Project 定义了项目的基本信息
type Project struct {
	BaseModel

	// 基本信息
	Title       string `gorm:"unique;size:128;not null"`
	Description string `gorm:"type:text"`

	// 负责人/审核人信息
	PrincipalId PKey    `gorm:"index;not null"`
	FkPrincipal *Member `gorm:"foreignKey:PrincipalId"`

	// 作品集相关信息
	WorksetId    PKey     `gorm:"not null;index"`
	FkWorkset    *Workset `gorm:"foreignKey:WorksetId"`
	WorksetIndex uint     `gorm:"not null"`
	// 旧版序号，保留对老作品的兼容
	LegacyId uint `gorm:"index"`

	// 龙译相关信息
	MoetranId string `gorm:"index;type:text"`

	// 当前项目的状态，全部拥有索引加速
	TranslateStatus uint8 `gorm:"not null;default:0"`     // 0: 未开始, 1: 翻译中, 2: 已翻译
	ProofStatus     uint8 `gorm:"not null;default:0"`     // 0: 未开始, 1: 校对中, 2: 已校对
	LetterStatus    uint8 `gorm:"not null;default:0"`     // 0: 未开始, 1: 字幕制作中, 2: 已完成
	ReviewStatus    uint8 `gorm:"not null;default:0"`     // 0: 未开始, 1: 审核中, 2: 已审核
	IsPublished     bool  `gorm:"not null;default:false"` // 是否已发布

	// 是否允许自动加入
	AllowAutoJoin bool `gorm:"not null;default:false"`
	// 是否为隐藏项目
	IsHidden bool `gorm:"not null;default:false"`
}

// 外部调用用的指针，避免每一次都创建
var emptyProject = &Project{}

// GetProject 获取一个空的 Project 实例
func GetProject() *Project {
	return emptyProject
}

// TableName 返回 Project 的表名
func (*Project) TableName() string {
	return "projects"
}

// ProjectSpec 定义了项目的查询条件
type ProjectSpec struct {
	Id *PKey

	Title *string // Title 只用于 LIKE 模糊查找，不在 map 中加入

	PrincipalId *PKey

	WorksetId    *PKey
	WorksetIndex *uint

	MoetranId *string

	TranslateStatus *[]uint8
	ProofreadStatus *[]uint8
	LetterStatus    *[]uint8
	ReviewStatus    *[]uint8
	IsPublished     *[]bool

	AllowAutoJoin *bool
	IsHidden      *bool
}

// Apply 将 ProjectSpec 应用到 WHERE 子句
func (s *ProjectSpec) Apply(query *gorm.DB) error {
	condition := make(map[string]any)

	if s.Title != nil {
		// Title 使用 LIKE 模糊查询
		query = query.Where("title LIKE ?", "%"+*s.Title+"%")
	}

	if s.Id != nil {
		condition["id"] = *s.Id
	}

	if s.PrincipalId != nil {
		condition["principal_id"] = *s.PrincipalId
	}

	// 作品集相关条件
	if s.WorksetId != nil {
		condition["workset_id"] = *s.WorksetId
	}
	if s.WorksetIndex != nil {
		condition["workset_index"] = *s.WorksetIndex
	}

	if s.MoetranId != nil {
		condition["moetran_id"] = *s.MoetranId
	}

	// 状态条件
	if s.TranslateStatus != nil {
		condition["translate_status"] = *s.TranslateStatus
	}
	if s.ProofreadStatus != nil {
		condition["proof_status"] = *s.ProofreadStatus
	}
	if s.LetterStatus != nil {
		condition["letter_status"] = *s.LetterStatus
	}
	if s.ReviewStatus != nil {
		condition["review_status"] = *s.ReviewStatus
	}
	if s.IsPublished != nil {
		condition["is_published"] = *s.IsPublished
	}

	// 项目权限条件
	if s.AllowAutoJoin != nil {
		condition["allow_auto_join"] = *s.AllowAutoJoin
	}
	if s.IsHidden != nil {
		condition["is_hidden"] = *s.IsHidden
	}

	return nil
}

// ProjectFields 定义了项目的字段，用于查询时选择特定字段返回
type ProjectFields struct {
	// BaseModel 字段
	Id   bool
	Time bool

	// Project 字段
	Title       bool
	Description bool

	PrincipalId bool
	UserFields  *UserFields // 需要 Preload 时指定

	MoetranId bool

	LegacyId     bool
	WorksetId    bool
	WorksetIndex bool

	Status bool // 包含所有状态

	AllowAutoJoin bool
	IsHidden      bool
}

// Apply 将 ProjectFields 转换为字符串切片并应用到 SELECT 子句
func (p *ProjectFields) Apply(query *gorm.DB) error {
	fields := []string{}

	if p.Id {
		fields = append(fields, "id")
	}
	if p.Time {
		fields = append(fields, "created_at", "updated_at", "deleted_at")
	}

	if p.Title {
		fields = append(fields, "title")
	}
	if p.Description {
		fields = append(fields, "description")
	}

	// 对负责人外键进行特殊处理
	if p.PrincipalId {
		fields = append(fields, "principal_id")
		if p.UserFields != nil {
			// 如果指定了 UserFields，则添加用户相关字段
			query = query.Preload("FkPrincipal.FkUser", func(db *gorm.DB) *gorm.DB {
				p.UserFields.Apply(db)
				return db
			})
		}
	}

	if p.MoetranId {
		fields = append(fields, "moetran_id")
	}

	if p.WorksetId {
		fields = append(fields, "workset_id")
		// 直接 Preload Workset 关联
		query = query.Preload("FkWorkset")
	}
	if p.WorksetIndex {
		fields = append(fields, "workset_index")
	}
	if p.LegacyId {
		fields = append(fields, "legacy_id")
	}

	if p.Status {
		fields = append(fields,
			"translate_status", "proofread_status",
			"letter_status", "review_status",
			"is_published")
	}

	if p.AllowAutoJoin {
		fields = append(fields, "allow_auto_join")
	}
	if p.IsHidden {
		fields = append(fields, "is_hidden")
	}

	return nil
}

// Insert 创建一个新的 Project 实例
func (*Project) Insert(hdl *gorm.DB, project *Project) error {
	if project == nil {
		return &InvalidParameterError{}
	}

	if project.Title == "" || project.WorksetId == 0 {
		return &LackOfRequiredFieldError{}
	}

	return hdl.
		Model(&Project{}).
		Create(project).Error
}

// Update 更新 Project 实例
func (*Project) Update(hdl *gorm.DB, project *Project) error {
	if project == nil {
		return &InvalidParameterError{}
	}

	if project.BaseModel.Id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Project{}).
		Where("id = ?", project.BaseModel.Id).
		Updates(project).Error
}

// SelectFirst 查询一个 Project 实例
func (*Project) SelectFirst(
	hdl *gorm.DB, cnd *ProjectSpec, fields *ProjectFields,
) (*Project, error) {
	var project Project

	query := hdl.Model(&Project{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if err := query.
		First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

// SelectMany 查询多个 Project 实例
func (*Project) SelectMany(
	hdl *gorm.DB, cnd *ProjectSpec, fields *ProjectFields,
	offset *int, limit *int,
) ([]*Project, error) {
	var projects []*Project

	query := hdl.Model(&Project{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.
		Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// Delete 删除一个 Project 实例
func (*Project) Delete(hdl *gorm.DB, id PKey) error {
	if id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Project{}).
		Where("id = ?", id).
		Delete(&Project{}).Error
}
