package repos

import (
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// AppliRepo 定义了申请相关的仓库接口
type AppliRepo interface {
	// GetApplicationSentByUserId 根据用户 ID 获取其发送的申请信息
	GetApplicationSentByUserId(userId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.ProjectApplication, error)
	// GetApplicationRecievedByUserId 根据用户 ID 获取收到的申请信息
	// 其原理是根据 principal_id 字段查询 labor_divisions 表，将其负责项目上的邀请信息查询出来
	GetApplicationRecievedByUserId(userId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.ProjectApplication, error)

	// CreateApplication 创建一个新的申请
	CreateApplication(applicant dbmodels.PrimaryKey, projectId dbmodels.PrimaryKey, targetRole uint) error
}

// appliRepoImpl 是 AppliRepo 的实现
type appliRepoImpl struct {
	handle *gorm.DB
}

// NewAppliRepo 创建一个新的 AppliRepo 实例
func NewAppliRepo(handle *gorm.DB) AppliRepo {
	return &appliRepoImpl{
		handle: handle,
	}
}

// GetApplicationSentByUserId 实现 AppliRepo 接口的 GetApplicationSentByUserId 方法
func (r *appliRepoImpl) GetApplicationSentByUserId(userId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.ProjectApplication, error) {
	var applications []*dbmodels.ProjectApplication

	if err := r.handle.Model(&dbmodels.ProjectApplication{}).
		Where("applicant_id = ?", userId).
		// 预加载 FkProject 关联的 Project 信息
		Preload("FkProject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "workset_id", "workset_index", "status")
		}).
		Find(&applications).
		Error; err != nil {
		return nil, err
	}

	return applications, nil
}

// GetApplicationRecievedByUserId 实现 AppliRepo 接口的 GetApplicationRecievedByUserId 方法
func (r *appliRepoImpl) GetApplicationRecievedByUserId(userId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.ProjectApplication, error) {
	var applications []*dbmodels.ProjectApplication

	// 查找 userId 作为 creator 的项目对应的申请
	if err := r.handle.Model(&dbmodels.ProjectApplication{}).
		Where("project_id IN (SELECT project_id FROM project_labor_divisions WHERE principal_id = ?)", userId).
		// 预加载 FkProject 关联的 Project 信息
		Preload("FkProject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "workset_id", "workset_index", "status")
		}).
		Offset(offset).
		Limit(limit).
		Find(&applications).
		Error; err != nil {
		return nil, err
	}

	return applications, nil
}

// CreateApplication 实现 AppliRepo 接口的 CreateApplication 方法
func (r *appliRepoImpl) CreateApplication(applicant dbmodels.PrimaryKey, projectId dbmodels.PrimaryKey, targetRole uint) error {
	application := &dbmodels.ProjectApplication{
		ApplicantId: applicant,
		ProjectId:   projectId,
		TargetRole:  dbmodels.LaborMask(targetRole),
	}

	if err := r.handle.Model(&dbmodels.ProjectApplication{}).
		Create(application).
		Error; err != nil {
		return err
	}

	return nil
}
