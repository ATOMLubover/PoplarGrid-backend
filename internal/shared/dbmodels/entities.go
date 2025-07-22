package dbmodels

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 成员的基本信息
type User struct {
	BaseModel

	// 基本信息
	Nickname     string `gorm:"uniqueIndex;size:128;not null"`
	Email        string `gorm:"unique;size:128;not null"`
	PasswordHash string `gorm:"size:256;not null"`

	// 尨译分配的 ID
	MoetranId string `gorm:"uniqueIndex;type:text"`
	// 龙译的 JWT Token
	MoetranAuth string `gorm:"type:text"`

	// 在仪表盘中的身份
	PoplarIsAdmin bool

	// 补充备注
	Remark string `gorm:"type:text"`
	// QQ 号
	QqNumber int

	// 上一次活跃时间（可能是通过 ping 来确定）
	LastActive time.Time
}

func (User) TableName() string {
	return "users"
}

// 汉化组
// 数量较少，可以一次加载完全，所以同样没有对 name 使用索引
type Team struct {
	BaseModel

	Name      string `gorm:"unique;size:256;not null"`
	MoetranId string `gorm:"uniqueIndex;type:text"`
}

func (Team) TableName() string {
	return "teams"
}

// 作品集（同步尨译信息）
// 考虑到数量较少，前端可以直接加载完整列表，这里不提供除了
// 主键和外键以外的索引
type Workset struct {
	BaseModel

	TeamId PrimaryKey
	FkTeam Team `gorm:"foreignKey:TeamId"`

	Name      string `gorm:"unique;type:text;not null"`
	MoetranId string `gorm:"uniqueIndex;type:text"`

	// 用于项目的组内自增序列
	ProjectSequenceName string `gorm:"uniqueIndex;type:text;not null"`
}

func (Workset) TableName() string {
	return "worksets"
}

// const PROJ_IDX_SEQ_PREFIX_FMT = "workset_project_index_seq_%d"

// // AfterCreate 在 workset 创建成功后为其创建一个 index 序列
// func (w *Workset) AfterCreate(tx *gorm.DB) (err error) {
// 	// 生成唯一的序列名称
// 	sequenceName := fmt.Sprintf(PROJ_IDX_SEQ_PREFIX_FMT, w.Id)

// 	// 动态执行 SQL 创建新的 PostgreSQL 序列
// 	createSequenceSQL := fmt.Sprintf("CREATE SEQUENCE %s START 1", sequenceName)
// 	if err := tx.Exec(createSequenceSQL).Error; err != nil {
// 		// 如果创建序列失败，返回错误，回滚整个 create 事务
// 		return fmt.Errorf("为 workset %d 生成序列 %s 失败: %w", w.Id, sequenceName, err)
// 	}

// 	// 将生成的序列名称赋值给 Workset 实例的 ProjectSequenceName 字段
// 	w.ProjectSequenceName = sequenceName
// 	if err := tx.Save(w).Error; err != nil {
// 		return fmt.Errorf("无法更新 workset %d 的序列记录名 %s: %w",
// 			w.Id, sequenceName, err)
// 	}

// 	return nil
// }

// 项目进度表
type Project struct {
	BaseModel

	// 内嵌作品的信息
	Title       string `gorm:"index;type:text;not null"`
	Description string `gorm:"type:text"`
	MoetranId   string `gorm:"uniqueIndex;type:text"`

	// 历史遗留序号（【】中的序号），保留对老作品的兼容
	// 经过观察，有序号重复的地方，如果可以最好重构这部分
	LegacyId uint `gorm:"index"`

	// 所属作品集
	WorksetId PrimaryKey
	FkWorkset Workset `gorm:"foreignKey:WorksetId"`

	// 项目在 workset 组内自增序列
	WorksetIndex uint `gorm:"not null"`

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

	// 负责人/审核人
	PrincipalId PrimaryKey `gorm:"index;not null"`
	FkPrincipal User       `gorm:"foreignKey:PrincipalId"`
}

func (Project) TableName() string {
	return "projects"
}

// // BeforeCreate 是在插入新项目之前触发的 hook
// // 其本质是为了动态地获取一个 project 在其 workset 中的唯一序列号
// func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
// 	// 根据当前新插入 p 的 workset_id 在 worksets 表查找到对应的序列名称
// 	var workset Workset

// 	if err := tx.Model(&Workset{}).
// 		Select("project_sequence_name").
// 		Where("id = ?", p.WorksetId).
// 		First(&workset).
// 		Error; err != nil {
// 		if err == gorm.ErrRecordNotFound || workset.ProjectSequenceName == "" {
// 			return fmt.Errorf("未找到 workset_id %d 对应的 sequence 名", p.WorksetId)
// 		}
// 		return fmt.Errorf("查找 workset_id %d 的 sequence 名时错误: %w", p.WorksetId, err)
// 	}

// 	// 使用 Raw 执行动态的 SQL 查询来获取下一个序列值
// 	var nextVal int64

// 	sqlQuery := fmt.Sprintf("SELECT nextval('%s')", workset.ProjectSequenceName)
// 	if err := tx.Raw(sqlQuery).Scan(&nextVal).Error; err != nil {
// 		return fmt.Errorf("无法从序列 %s 中获取 nextval: %w", workset.ProjectSequenceName, err)
// 	}

// 	// 将获取到的下一个序列值赋值给新插入项目记录的 workset_index
// 	p.WorksetIndex = uint(nextVal)

// 	// 继续执行 create 流
// 	return nil
// }

// AfterDelete 在软删除项目后触发，是为了级联（软）删除对应表
func (p *Project) AfterDelete(tx *gorm.DB) (err error) {
	// 删除对应的 labor_division 记录
	if err := tx.Where("project_id = ?", p.Id).
		Delete(&ProjectLaborDivision{}).Error; err != nil {
		return fmt.Errorf("删除 project_id %d 的 labor_division 记录失败: %w", p.Id, err)
	}

	// 删除对应的 invitation 记录
	if err := tx.Where("project_id = ?", p.Id).
		Delete(&ProjectInvitation{}).Error; err != nil {
		return fmt.Errorf("删除 project_id %d 的 invitation 记录失败: %w", p.Id, err)
	}

	// 删除对应的 application 记录
	if err := tx.Where("project_id = ?", p.Id).
		Delete(&ProjectApplication{}).Error; err != nil {
		return fmt.Errorf("删除 project_id %d 的 application 记录失败: %w", p.Id, err)
	}

	return nil
}
