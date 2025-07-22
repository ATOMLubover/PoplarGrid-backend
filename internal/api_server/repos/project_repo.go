package repos

import (
	"errors"
	"poplargrid/internal/api_server/dtos"
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// project 所需基础字段
var kProjectBasicFields = []string{
	"projects.id",
	"projects.title",

	"projects.legacy_id",
	"projects.moetran_id",
	"projects.workset_id",
	"projects.workset_index",

	"projects.principal_id",

	"projects.translate_status",
	"projects.proof_status",
	"projects.letter_status",
	"projects.review_status",
	"projects.is_published",

	"projects.allow_auto_join",
	"projects.created_at",
	"projects.updated_at",
}

// ProjectRepo 接口定义了项目仓库的基本操作
type ProjectRepo interface {
	Repo

	// SelectBasicPageWithParam 按 ID 倒序获取项目列表，支持分页以及复合条件查询
	SelectBasicPageWithParam(offset, limit int, queryParams *dtos.ProjectSearchParams) ([]*dbmodels.Project, error)

	// SelectById 获取指定 ID 的项目的全部信息
	SelectById(id dbmodels.PrimaryKey) (*dbmodels.Project, error)

	// CreateProject 创建一个新的项目，如果成功则 project 参数的 Id、WorksetIndex 字段会被填充
	// 同时 workset 和相关的 team 信息会被递归加载
	CreateProject(project *dbmodels.Project) error

	// UpdateMoetranId 更新项目的 MoetranId 字段
	UpdateMoetranId(projectId dbmodels.PrimaryKey, moetranId string) error
	// SaveInfo 更新项目的基本信息
	SaveInfo(project *dbmodels.Project) error

	// DeleteById 删除指定 ID 的项目
	DeleteById(id dbmodels.PrimaryKey) error
}

// projectRepoImpl 是 ProjectRepo 的实现
type projectRepoImpl struct {
	handle *gorm.DB
}

// NewProjectRepo 创建一个新的 ProjectRepo 实例
func NewProjectRepo(
	db *gorm.DB,
) ProjectRepo {
	return &projectRepoImpl{
		handle: db,
	}
}

// Table 限定当前操作的表名
func (r *projectRepoImpl) Table() *gorm.DB {
	return r.handle.Table(dbmodels.Project{}.TableName())
}

// GetHandle 实现 ProjectRepo.Repo 接口的 GetHandle 方法
func (r *projectRepoImpl) GetHandle() *gorm.DB {
	return r.handle
}

// SelectBasicPageWithParam 实现 ProjectRepo 接口的 SelectBasicPageIdDescWithParam 方法
func (r *projectRepoImpl) SelectBasicPageWithParam(offset, limit int, queryParams *dtos.ProjectSearchParams) ([]*dbmodels.Project, error) {
	var projects []*dbmodels.Project

	// 构建查询条件
	query := r.buildQueryWithParams(r.Table(), queryParams)

	// 最后进行查询
	if err := query.
		Select(kProjectBasicFields).
		Limit(limit).
		Offset(offset).
		Find(&projects).
		Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// SelectById 实现 ProjectRepo 接口的 SelectById 方法
func (r *projectRepoImpl) SelectById(id dbmodels.PrimaryKey) (*dbmodels.Project, error) {
	var project dbmodels.Project

	if err := r.Table().
		Where("id = ?", id).
		First(&project).
		Error; err != nil {
		return nil, err
	}

	return &project, nil
}

// CreateProject 实现 ProjectRepo 接口的 CreateProject 方法
func (r *projectRepoImpl) CreateProject(project *dbmodels.Project) error {
	if project == nil {
		return errors.New("project 不能为 nil")
	}

	// 先尝试插入新项目，这会导致相关的 trigger 和 function 被触发
	if err := r.Table().Create(project).Error; err != nil {
		return err
	}

	// 随后尝试获取对应的 workset_index
	// 递归加载 workset 信息和 team 信息
	if err := r.Table().
		Preload("FkWorkset", func(db *gorm.DB) {
			db.
				Preload("FkTeam", func(db *gorm.DB) {
					db.Select("id", "moetran_id")
				}).
				Select("id", "moetran_id")
		}).
		Where("id = ?", project.Id).
		Find(project).
		Error; err != nil {
		return err
	}

	// 在这里 workset_index 应当被填充
	return nil
}

// UpdateMoetranId 实现 ProjectRepo 接口的 UpdateMoetranId 方法
func (r *projectRepoImpl) UpdateMoetranId(projectId dbmodels.PrimaryKey, moetranId string) error {
	// 执行更新操作
	if err := r.Table().
		Where("id = ?", projectId).
		Update("moetran_id", moetranId).
		Error; err != nil {
		return err
	}

	return nil
}

// SaveInfo 实现 ProjectRepo 接口的 SaveInfo 方法
func (r *projectRepoImpl) SaveInfo(project *dbmodels.Project) error {
	if project == nil {
		return errors.New("project 不能为 nil")
	}

	// 执行更新操作
	if err := r.Table().
		// 这里使用 Save 方法会自动处理主键和更新字段，利用零值保护简化
		Save(project).
		Error; err != nil {
		return err
	}

	return nil
}

// DeleteById 实现 ProjectRepo 接口的 DeleteById 方法
func (r *projectRepoImpl) DeleteById(id dbmodels.PrimaryKey) error {
	// 执行删除操作
	if err := r.Table().
		Where("id = ?", id).
		// 这里强制使用软删除
		Delete(&dbmodels.Project{}).
		Error; err != nil {
		return err
	}

	return nil
}

// ================ 辅助函数 ================

// buildQueryWithParams 通过查询参数 query 构建 WHERE 子句
func (r *projectRepoImpl) buildQueryWithParams(base *gorm.DB, queryParams *dtos.ProjectSearchParams) *gorm.DB {
	if queryParams == nil {
		return base
	}

	// 先确定是否需要根据 user 查询
	if queryParams.UserId != nil {
		// 如果需要根据 user 查询，则需要根据 labor 表使用 inner join
		base = base.
			Joins("INNER JOIN project_labor_divisions ON project_labor_divisions.project_id = projects.id").
			Where("project_labor_divisions.user_id = ?", queryParams.UserId)
	}

	// 如果不需要根据 user 查询，则进行普通查询
	// 先添加工作集 ID 的查询条件
	if queryParams.WorksetId != nil {
		base = base.Where("workset_id = ?", *queryParams.WorksetId)
	}

	// 添加 sort 条件
	switch queryParams.Sort {
	case dtos.SORT_ID_DESC:
		base = base.Order("id DESC")
	case dtos.SORT_UPDATED_AT_DESC:
		base = base.Order("updated_at DESC")
	default:
		// 默认按 ID 倒序
		base = base.Order("id DESC")
	}

	// 添加状态查询条件
	if queryParams.Status != nil {
		if queryParams.Status.TranslateStatus != nil {
			base = base.Where("translate_status = ?", *queryParams.Status.TranslateStatus)
		}
		if queryParams.Status.ProofStatus != nil {
			base = base.Where("proof_status = ?", *queryParams.Status.ProofStatus)
		}
		if queryParams.Status.LetterStatus != nil {
			base = base.Where("letter_status = ?", *queryParams.Status.LetterStatus)
		}
		if queryParams.Status.ReviewStatus != nil {
			base = base.Where("review_status = ?", *queryParams.Status.ReviewStatus)
		}
		if queryParams.Status.PublishStatus != nil {
			base = base.Where("is_published = ?", *queryParams.Status.PublishStatus)
		}
	}

	return base
}
