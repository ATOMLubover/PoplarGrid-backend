package models

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

// NewLaborMask 根据多个职责掩码创建一个新的复合掩码
func NewLaborMask(roles ...LaborMask) LaborMask {
	var mask LaborMask

	for _, role := range roles {
		mask.AddRole(role)
	}

	return mask
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

//
