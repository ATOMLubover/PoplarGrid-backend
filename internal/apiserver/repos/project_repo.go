package repos

import (
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// project 所需基础字段
var kProjectBasicFields = []string{
	"id",
	"title",
	"legacy_id",
	"workset_id",
	"workset_index",
	"translate_status",
	"proof_status",
	"letter_status",
	"review_status",
	"is_published",
	"allow_auto_join",
	"created_at",
	"updated_at",
}

// ProjectRepo 接口定义了项目仓库的基本操作
type ProjectRepo interface {
	Repo

	// SelectBasicPageIdDescWithParam 按 ID 倒序获取项目列表，支持分页以及复合条件查询
	SelectBasicPageIdDescWithParam(worksetId dbmodels.PrimaryKey, offset, limit int, queryParams *dtos.ProjectStatusQueryParams) ([]*dbmodels.Project, error)
	// SelectBasicPageUpdatedAtDescWithParam 按更新时间倒序获取项目列表，支持分页以及复合条件查询
	SelectBasicPageUpdatedAtDescWithParam(worksetId dbmodels.PrimaryKey, offset, limit int, queryParams *dtos.ProjectStatusQueryParams) ([]*dbmodels.Project, error)

	// SelectById 获取指定 ID 的项目的全部信息
	SelectById(id dbmodels.PrimaryKey) (*dbmodels.Project, error)

	// CreateProject 创建一个新的项目，如果成功则 project 参数的 Id、WorksetIndex 字段会被填充
	CreateProject(project *dbmodels.Project) error
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

// SelectBasicPageIdDescWithParam 实现 ProjectRepo 接口的 SelectBasicPageIdDescWithParam 方法
func (r *projectRepoImpl) SelectBasicPageIdDescWithParam(worksetId dbmodels.PrimaryKey, offset, limit int, queryParams *dtos.ProjectStatusQueryParams) ([]*dbmodels.Project, error) {
	var projects []*dbmodels.Project

	// 构建查询条件
	query := r.buildQueryWithParams(r.Table(), queryParams)

	// 最后进行查询
	if err := query.
		Where("workset_id = ?", worksetId).
		Select(kProjectBasicFields).
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).
		Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// SelectBasicPageUpdatedAtDescWithParam 实现 ProjectRepo 接口的 SelectBasicPageUpdatedAtDescWithParam 方法
func (r *projectRepoImpl) SelectBasicPageUpdatedAtDescWithParam(worksetId dbmodels.PrimaryKey, offset, limit int, queryParams *dtos.ProjectStatusQueryParams) ([]*dbmodels.Project, error) {
	var projects []*dbmodels.Project

	// 构建查询条件
	query := r.buildQueryWithParams(r.Table(), queryParams)

	// 最后进行查询
	if err := query.
		Where("workset_id = ?", worksetId).
		Select(kProjectBasicFields).
		Order("updated_at DESC").
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

// ================ 辅助函数 ================

// buildQueryWithParams 通过查询参数为 query 构建 WHERE 子句
func (r *projectRepoImpl) buildQueryWithParams(base *gorm.DB, queryParams *dtos.ProjectStatusQueryParams) *gorm.DB {
	if queryParams == nil {
		return base
	}

	if queryParams.TranslateStatus != nil {
		base = base.Where("translate_status = ?", *queryParams.TranslateStatus)
	}
	if queryParams.ProofStatus != nil {
		base = base.Where("proof_status = ?", *queryParams.ProofStatus)
	}
	if queryParams.LetterStatus != nil {
		base = base.Where("letter_status = ?", *queryParams.LetterStatus)
	}
	if queryParams.ReviewStatus != nil {
		base = base.Where("review_status = ?", *queryParams.ReviewStatus)
	}
	if queryParams.PublishStatus != nil {
		base = base.Where("is_published = ?", *queryParams.PublishStatus)
	}

	return base
}

// CreateProject 实现 ProjectRepo 接口的 CreateProject 方法
func (r *projectRepoImpl) CreateProject(project *dbmodels.Project) error {
	// 先尝试插入新项目，这会导致相关的 trigger 和 function 被触发
	if err := r.Table().Create(project).Error; err != nil {
		return err
	}

	// 随后尝试获取对应的 workset_index
	if err := r.Table().
		Select("id", "workset_index").
		Where("id = ?", project.Id).
		Find(project).
		Error; err != nil {
		return err
	}

	// 在这里 workset_index 应当被填充
	return nil
}
