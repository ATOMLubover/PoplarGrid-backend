package models

import "gorm.io/gorm"

// Application 定义了申请的基本信息
type Application struct {
	BaseModel

	ApplicantMemberId PKey    `gorm:"not null;index"`               // 申请者成员 ID
	FkApplicant       *Member `gorm:"foreignKey:ApplicantMemberId"` // 申请者成员外键

	ProcessorMemberId PKey    `gorm:"not null;index"`               // 负责人成员 ID
	FkProcessor       *Member `gorm:"foreignKey:ProcessorMemberId"` // 负责人成员外键

	TargetProjectId PKey     `gorm:"not null;index"`             // 目标项目 ID
	FkProject       *Project `gorm:"foreignKey:TargetProjectId"` // 目标项目外键

	TargetLaborMask uint32 `gorm:"not null;default:0"` // 目标角色掩码

	Status Status `gorm:"not null;default:0"` // 申请状态，0: 待处理, 1: 已接受, 2: 已拒绝
}

// TableName 返回 Application 的表名
func (*Application) TableName() string {
	return "applications"
}

// ApplicationSpec 定义了申请的查询条件
type ApplicationSpec struct {
	Id *PKey // 申请 ID

	ApplicantMemberId *PKey // 申请者成员 ID
	ProcessorMemberId *PKey // 负责人成员 ID

	TargetProjectId *PKey // 目标项目 ID

	Status *Status // 申请状态
}

// Apply 将 ApplicationSpec 应用为 WHERE 子句
func (s *ApplicationSpec) Apply(query *gorm.DB) {
	if s.Id != nil {
		query = query.Where("id = ?", *s.Id)
	}

	if s.ApplicantMemberId != nil {
		query = query.Where("applicant_member_id = ?", *s.ApplicantMemberId)
	}
	if s.ProcessorMemberId != nil {
		query = query.Where("processor_member_id = ?", *s.ProcessorMemberId)
	}

	if s.TargetProjectId != nil {
		query = query.Where("target_project_id = ?", *s.TargetProjectId)
	}

	if s.Status != nil {
		query = query.Where("status = ?", *s.Status)
	}
}

// ApplicationFields 定义了 Application 的预加载字段
type ApplicationFields struct {
	Id   bool // 申请 ID
	Time bool // 申请时间

	ApplicantMemberId bool          // 申请者成员 ID
	ProcessorMemberId bool          // 负责人成员 ID
	MemberFields      *MemberFields // 申请者成员信息

	TargetProjectId bool           // 目标项目 ID
	ProjectFields   *ProjectFields // 目标项目信息

	TargetLaborMask bool // 目标角色掩码

	Status bool // 申请状态
}

// Apply 将 ApplicationFields 应用为 SELECT 子句
func (f *ApplicationFields) Apply(query *gorm.DB) {
	fields := []string{}

	if f.Id {
		fields = append(fields, "id")
	}
	if f.Time {
		fields = append(fields, "created_at", "updated_at", "deleted_at")
	}

	if f.ApplicantMemberId {
		fields = append(fields, "applicant_member_id")
		if f.MemberFields != nil {
			query = query.Preload("FkApplicant.FkUser", func(db *gorm.DB) *gorm.DB {
				f.MemberFields.Apply(db)
				return db
			})
		}
	}
	if f.ProcessorMemberId {
		fields = append(fields, "processor_member_id")
		if f.MemberFields != nil {
			query = query.Preload("FkProcessor.FkUser", func(db *gorm.DB) *gorm.DB {
				f.MemberFields.Apply(db)
				return db
			})
		}
	}

	if f.TargetProjectId {
		fields = append(fields, "target_project_id")
		if f.ProjectFields != nil {
			query = query.Preload("FkProject", func(db *gorm.DB) *gorm.DB {
				f.ProjectFields.Apply(db)
				return db
			})
		}
	}

	if f.TargetLaborMask {
		fields = append(fields, "target_labor_mask")
	}

	if f.Status {
		fields = append(fields, "status")
	}
}

// 供外部调用的指针，避免每一次都创建
var emptyApplication *Application = nil

// GetApplication 获取一个空的 Application 实例
func GetApplication() *Application {
	return emptyApplication
}

// Insert 插入 Application 实例到数据库
func (*Application) Insert(hdl *gorm.DB, application *Application) error {
	if application == nil {
		return &InvalidParameterError{}
	}

	return hdl.
		Model(&Application{}).
		Create(application).Error
}

// Update 更新 Application 实例到数据库
func (*Application) Update(hdl *gorm.DB, application *Application) error {
	if application == nil {
		return &InvalidParameterError{}
	}

	if application.BaseModel.Id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Application{}).
		Where("id = ?", application.BaseModel.Id).
		Updates(application).Error
}

// SelectFirst 查询一个 Application 实例
func (*Application) SelectFirst(
	hdl *gorm.DB, cnd *ApplicationSpec, fields *ApplicationFields,
) (*Application, error) {
	var application Application

	query := hdl.Model(&Application{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if err := query.
		First(&application).Error; err != nil {
		return nil, err
	}

	return &application, nil
}

// SelectMany 查询多个 Application 实例
func (*Application) SelectMany(
	hdl *gorm.DB, cnd *ApplicationSpec, fields *ApplicationFields,
) ([]*Application, error) {
	var applications []*Application

	query := hdl.Model(&Application{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if err := query.
		Find(&applications).Error; err != nil {
		return nil, err
	}

	return applications, nil
}

// Delete 删除 Application 实例
func (*Application) Delete(hdl *gorm.DB, id PKey) error {
	if id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&Application{}).
		Where("id = ?", id).
		Delete(&Application{}).Error
}
