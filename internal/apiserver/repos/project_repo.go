package repos

import (
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

const kProjectBasicFields = "id, " +
	"title, " +
	"legacy_id, " +
	"workset_id, " +
	"workset_index, " +
	"translate_status, " +
	"proof_status, " +
	"letter_status, " +
	"review_status, " +
	"is_published" +
	"allow_auto_join, " +
	"created_at, " +
	"updated_at"

// ProjectRepo 接口定义了项目仓库的基本操作
type ProjectRepo interface {
	// SelectBasicPage 按 ID 倒序获取所有项目列表，包括已发布和未发布的项目
	SelectBasicPageIdDesc(worksetId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.Project, error)
	// SelectBasicPageUpdatedAtDesc 按更新时间倒序获取项目列表，包括已发布和未发布的项目
	SelectBasicPageUpdatedAtDesc(worksetId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.Project, error)
	// SelectBasicPageIdDescWithParam 按 ID 倒序获取项目列表，支持分页和排序以及复合条件查询
	SelectBasicPageIdDescWithParam(worksetId dbmodels.PrimaryKey, offset, limit int, queryParams *dtos.ProjectStatusQueryParams) ([]*dbmodels.Project, error)
	// SelectBasicPageUpdatedAtDescWithParam 按更新时间倒序获取项目列表，支持分页和排序以及复合条件查询
	SelectBasicPageUpdatedAtDescWithParam(worksetId dbmodels.PrimaryKey, offset, limit int, queryParams *dtos.ProjectStatusQueryParams) ([]*dbmodels.Project, error)
}

// projectRepoImpl 是 ProjectRepo 的实现
type projectRepoImpl struct {
	db *gorm.DB
}

// NewProjectRepo 创建一个新的 ProjectRepo 实例
func NewProjectRepo(
	db *gorm.DB,
) ProjectRepo {
	return &projectRepoImpl{
		db: db,
	}
}

// Table 限定当前操作的表名
func (r *projectRepoImpl) Table() *gorm.DB {
	return r.db.Table(dbmodels.Project{}.TableName())
}

// SelectBasicPageIdDesc 实现 ProjectRepo 接口的 SelectBasicPageIdDesc 方法
func (r *projectRepoImpl) SelectBasicPageIdDesc(worksetId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.Project, error) {
	var projects []*dbmodels.Project

	if err := r.Table().
		Where("workset_id = ?", worksetId).
		Select(kProjectBasicFields).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&projects).
		Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// SelectBasicPageUpdatedAtDesc 实现 ProjectRepo 接口的 SelectBasicPageUpdatedAtDesc 方法
func (r *projectRepoImpl) SelectBasicPageUpdatedAtDesc(worksetId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.Project, error) {
	var projects []*dbmodels.Project

	if err := r.Table().
		Where("workset_id = ?", worksetId).
		Select(kProjectBasicFields).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&projects).
		Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// SelectBasicPageIdDescWithParam 实现 ProjectRepo 接口的 SelectBasicPageIdDescWithParam 方法
func (r *projectRepoImpl) SelectBasicPageIdDescWithParam(worksetId dbmodels.PrimaryKey, offset, limit int, queryParams *dtos.ProjectStatusQueryParams) ([]*dbmodels.Project, error) {
	var projects []*dbmodels.Project

	// 逐步构建查询条件
	query := r.Table()
	if queryParams.TranslateStatus != nil {
		query = query.Where("translate_status = ?", *queryParams.TranslateStatus)
	}
	if queryParams.ProofStatus != nil {
		query = query.Where("proof_status = ?", *queryParams.ProofStatus)
	}
	if queryParams.LetterStatus != nil {
		query = query.Where("letter_status = ?", *queryParams.LetterStatus)
	}
	if queryParams.ReviewStatus != nil {
		query = query.Where("review_status = ?", *queryParams.ReviewStatus)
	}

	// 最后进行查询
	if err := query.
		Where("workset_id = ?", worksetId).
		Where("is_pulished = ?", false). // 只查找未发布的项目
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

	// 逐步构建查询条件
	query := r.Table()
	if queryParams.TranslateStatus != nil {
		query = query.Where("translate_status = ?", *queryParams.TranslateStatus)
	}
	if queryParams.ProofStatus != nil {
		query = query.Where("proof_status = ?", *queryParams.ProofStatus)
	}
	if queryParams.LetterStatus != nil {
		query = query.Where("letter_status = ?", *queryParams.LetterStatus)
	}
	if queryParams.ReviewStatus != nil {
		query = query.Where("review_status = ?", *queryParams.ReviewStatus)
	}

	// 最后进行查询
	if err := query.
		Where("workset_id = ?", worksetId).
		Where("is_pulished = ?", false). // 只查找未发布的项目
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
