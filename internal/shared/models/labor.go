package models

import "gorm.io/gorm"

// LaborMask 定义了分工的掩码类型
type LaborMask uint32

const (
	LABOR_PRINCIPAL_MASK   LaborMask = 1 << iota // 创建者 + 负责人
	LABOR_SRC_PROV_MASK                          // 图源
	LABOR_PERFECTOR_MASK                         // 美工
	LABOR_TRANSLATOR_MASK                        // 翻译
	LABOR_PROOFREADER_MASK                       // 校对
	LABOR_LETTERER_MASK                          // 嵌字
	LABOR_REVIEWER_MASK                          // 嵌字审核
	LABOR_PUBLISHER_MASK                         // 发布者
)

// HasRole 检查 LaborMask 是否包含某个职责
func (m LaborMask) HasRole(roleMask LaborMask) bool {
	return (m & roleMask) != 0
}

// AddRole 为 LaborMask 添加一个职责
func (m *LaborMask) AddRole(roleMask LaborMask) {
	*m |= roleMask
}

// RemoveRole 从 LaborMask 中移除一个职责
func (m *LaborMask) RemoveRole(roleMask LaborMask) {
	*m &= ^roleMask
}

// Labor 定义了 project 中一个 member 的基本信息
type Labor struct {
	BaseModel

	// 分工，按掩码存储，因为没有查询需求
	LaborMask LaborMask `gorm:"not null;default:0"`

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
