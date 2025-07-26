package models

import "gorm.io/gorm"

// Member 定义了成员的基本信息
type Member struct {
	BaseModel

	// 组内职位
	IsAdmin          bool `gorm:"default:false"` // 是否为管理员
	IsSourceProvider bool `gorm:"default:false"` // 是否为图源
	IsPerfector      bool `gorm:"default:false"` // 是否为美工
	IsTranslator     bool `gorm:"default:false"` // 是否为翻译
	IsProofreader    bool `gorm:"default:false"` // 是否为校对
	IsLetterer       bool `gorm:"default:false"` // 是否为嵌字
	IsReviewer       bool `gorm:"default:false"` // 是否为嵌字审核
	IsPublisher      bool `gorm:"default:false"` // 是否为发布者

	// 实际用户外键
	UserId PKey  `gorm:"not null;index"`
	FkUser *User `gorm:"foreignKey:UserId"`

	// 所属团队外键
	TeamId PKey  `gorm:"not null;index"`
	FkTeam *Team `gorm:"foreignKey:TeamId"`
}

// 外部调用用的指针，避免每一次都创建
var emptyMember = &Member{}

// GetMember 获取一个空的 Member 实例
func GetMember() *Member {
	return emptyMember
}

// TableName 返回 Member 的表名
func (*Member) TableName() string {
	return "members"
}

// MemberSpec 定义了成员的查询条件
// 外键为 0 时，只 Preload 但是不施加条件
type MemberSpec struct {
	Id *PKey

	IsAdmin       *bool
	IsSource      *bool
	IsPerfector   *bool
	IsTranslator  *bool
	IsProofreader *bool
	IsLetterer    *bool
	IsReviewer    *bool
	IsPublisher   *bool

	UserId *PKey

	TeamId *PKey
}

// Apply 将 MemberSpec 应用为 WHERE 子句
func (s *MemberSpec) Apply(query *gorm.DB) {
	if s.Id != nil {
		query = query.Where("id = ?", *s.Id)
	}

	if s.IsAdmin != nil {
		query = query.Where("is_admin = ?", *s.IsAdmin)
	}
	if s.IsSource != nil {
		query = query.Where("is_source_provider = ?", *s.IsSource)
	}
	if s.IsPerfector != nil {
		query = query.Where("is_perfector = ?", *s.IsPerfector)
	}
	if s.IsTranslator != nil {
		query = query.Where("is_translator = ?", *s.IsTranslator)
	}
	if s.IsProofreader != nil {
		query = query.Where("is_proofreader = ?", *s.IsProofreader)
	}
	if s.IsLetterer != nil {
		query = query.Where("is_letterer = ?", *s.IsLetterer)
	}
	if s.IsReviewer != nil {
		query = query.Where("is_reviewer = ?", *s.IsReviewer)
	}
	if s.IsPublisher != nil {
		query = query.Where("is_publisher = ?", *s.IsPublisher)
	}

	if s.UserId != nil {
		query = query.Where("user_id = ?", *s.UserId)
	}

	if s.TeamId != nil {
		query = query.Where("team_id = ?", *s.TeamId)
	}
}

// MemberFields 定义了成员的字段，用于查询时选择特定字段返回
type MemberFields struct {
	// BaseModel 字段
	Id   bool
	Time bool

	// Member 字段
	IsAdmin bool
	Roles   bool // 该选项将返回所有职责相关字段

	// 用户外键
	UserId     bool
	UserFields *UserFields // 需要 Preload 时指定

	// 汉化组外键
	TeamId     bool
	TeamFields bool // 需要 Preload 时指定
}

// rolesFields 包含所有职责相关字段的映射
var rolesFields = []string{
	"is_source",
	"is_perfector",
	"is_translator",
	"is_proofreader",
	"is_letterer",
	"is_reviewer",
	"is_publisher",
}

// Apply 将 MemberFields 转换为字符串切片并应用为 SELECT 子句
func (m *MemberFields) Apply(query *gorm.DB) {
	fields := []string{}

	// BaseModel 字段
	if m.Id {
		fields = append(fields, "id")
	}
	if m.Time {
		fields = append(fields, "created_at", "updated_at", "deleted_at")
	}

	// Member 字段
	if m.IsAdmin {
		fields = append(fields, "is_admin")
	}
	if m.Roles {
		fields = append(fields, rolesFields...)
	}

	// 用户外键字段
	if m.UserId {
		fields = append(fields, "user_id")
		if m.UserFields != nil {
			// 如果指定了 UserFields，则添加用户相关字段
			query.Preload("FkUser", func(db *gorm.DB) *gorm.DB {
				m.UserFields.Apply(db)
				return db
			})
		}
		// 如果没有指定 UserFields，则只查询 UserId 字段，而不进行 Preload
	}
	// 汉化组外键字段
	if m.TeamId {
		fields = append(fields, "team_id")
		if m.TeamFields {
			// 如果指定了 TeamFields，则添加团队相关字段
			query.Preload("FkTeam")
		}
	}

	query = query.Select(fields)
}

// Insert 创建一个新的 Member 实例
func (*Member) Insert(hdl *gorm.DB, member *Member) error {
	if member == nil {
		return &InvalidParameterError{}
	}

	if member.FkUser == nil || member.FkTeam == nil {
		return &LackOfRequiredFieldError{}
	}

	return hdl.
		Model(&Member{}).
		Create(member).Error
}

// Update 更新 Member 实例
func (*Member) Update(hdl *gorm.DB, member *Member) error {
	if member == nil {
		return &InvalidParameterError{}
	}

	if member.BaseModel.Id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Member{}).
		Updates(member).Error
}

// SelectFirst 查询一个 Member 实例
func (*Member) SelectFirst(
	hdl *gorm.DB, cnd *MemberSpec, fields *MemberFields,
) (*Member, error) {
	var member Member

	query := hdl.Model(&Member{})

	if cnd != nil {
		// 如果提供了查询条件，则应用到查询中
		cnd.Apply(query)
	}

	if fields != nil {
		// 如果提供了字段选择，则应用到查询中
		fields.Apply(query)
	}

	if err := query.
		First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// SelectMany 查询多个 Member 实例
func (*Member) SelectMany(
	hdl *gorm.DB, cnd *MemberSpec, fields *MemberFields,
	offset *int, limit *int,
) ([]*Member, error) {
	var members []*Member

	query := hdl.Model(&Member{})

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
		Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

// Delete 删除一个 Member 实例
func (*Member) Delete(hdl *gorm.DB, id PKey) error {
	if id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Member{}).
		Where("id = ?", id).
		Delete(&Member{}).Error
}
