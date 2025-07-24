package repos

import (
	"fmt"
	"poplargrid/internal/api_server/dtos"
	"poplargrid/internal/shared/dbmodels"

	"github.com/kataras/iris/v12/x/errors"
	"gorm.io/gorm"
)

// MaterialView 接口定义了物化视图的基本操作
type MaterialView interface {
	// SelectProjectStats 获取项目统计视图
	SelectProjectStats(worksetId dbmodels.PrimaryKey) (*dbmodels.ProjectStats, error)
	// RefreshProjectStats 手动刷新项目统计视图，且支持并发刷新
	RefreshProjectStats() error
}

// materialViewImpl 是 MaterialView 的实现
type materialViewImpl struct {
	handle *gorm.DB
}

// NewMaterialView 创建一个新的 MaterialView 实例
func NewMaterialView(db *gorm.DB) MaterialView {
	return &materialViewImpl{
		handle: db,
	}
}

// SelectProjectStats 实现 MaterialView 接口的 SelectProjectStats 方法
func (m *materialViewImpl) SelectProjectStats(worksetId dbmodels.PrimaryKey) (*dbmodels.ProjectStats, error) {
	var stats *dbmodels.ProjectStats

	if err := m.handle.
		Table(dbmodels.ProjectStats{}.TableName()).
		Where("workset_id = ?", worksetId).
		Find(&stats).
		Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// RefreshProjectStats 实现 MaterialView 接口的 RefreshProjectStats 方法
func (m *materialViewImpl) RefreshProjectStats() error {
	// 使用原生 SQL 刷新物化视图
	if err := m.handle.
		Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY " +
			dbmodels.ProjectStats{}.TableName()).
		Error; err != nil {
		return err
	}

	return nil
}

// LaborRepo 接口定义了项目中的成员分工的仓库操作
type LaborRepo interface {
	// SelectByProjectId 获取指定项目 ID 的成员分工列表
	SelectByProjectId(projectId dbmodels.PrimaryKey) ([]*dbmodels.ProjectLaborDivision, error)
	// SelectByUserId 获取指定用户 ID 的成员分工列表
	// 这会预加载 FkUser 关联的 User 信息
	SelectByUserId(userId dbmodels.PrimaryKey, projectIds []dbmodels.PrimaryKey) ([]*dbmodels.ProjectLaborDivision, error)

	// // SelectProjectPageByUserId 获取用户参与的项目列表，按 ID 倒序，支持分页
	// SelectProjectPageByUserId(userId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.ProjectLaborDivision, error)

	// CreateLaborDivision 创建一个新的成员分工记录
	CreateLaborDivision(labor *dbmodels.ProjectLaborDivision) error
}

// laborRepoImpl 是 LaborRepo 的实现
type laborRepoImpl struct {
	handle *gorm.DB
}

// NewLaborRepo 创建一个新的 LaborRepo 实例
func NewLaborRepo(db *gorm.DB) LaborRepo {
	return &laborRepoImpl{
		handle: db,
	}
}

// SelectByProjectId 实现 LaborRepo 接口的 SelectByProjectId 方法
func (r *laborRepoImpl) SelectByProjectId(projectId dbmodels.PrimaryKey) ([]*dbmodels.ProjectLaborDivision, error) {
	var labors []*dbmodels.ProjectLaborDivision

	if err := r.handle.Model(&dbmodels.ProjectLaborDivision{}).
		// 这里使用 Preload 来加载用户信息
		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nickname") // User 只需要 id、nickname
		}).
		Where("project_id = ?", projectId).
		Find(&labors).
		Error; err != nil {
		return nil, err
	}

	return labors, nil
}

// SelectByUserId 实现 LaborRepo 接口的 SelectByUserId 方法
func (r *laborRepoImpl) SelectByUserId(userId dbmodels.PrimaryKey, projectIds []dbmodels.PrimaryKey) ([]*dbmodels.ProjectLaborDivision, error) {
	var labors []*dbmodels.ProjectLaborDivision

	if err := r.handle.Model(&dbmodels.ProjectLaborDivision{}).
		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nickname") // User 只需要 id、nickname
		}).
		Where("user_id = ? AND project_id IN (?)", userId, projectIds).
		Find(&labors).
		Error; err != nil {
		return nil, err
	}

	return labors, nil
}

// // SelectProjectPageByUserId 实现 LaborRepo 接口的 SelectProjectPageByUserId 方法
// func (r *laborRepoImpl) SelectProjectPageByUserId(userId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.ProjectLaborDivision, error) {
// 	var labors []*dbmodels.ProjectLaborDivision

// 	if err := r.handle.Model(&dbmodels.ProjectLaborDivision{}).
// 		// 预加载 FkProject 关联的 Project 信息
// 		Preload("FkProject", func(db *gorm.DB) *gorm.DB {
// 			return db.Select(kProjectBasicFields)
// 		}).
// 		Select("id", "project_id", "labor_role").
// 		Order("id DESC"). // 按 ID 倒序，也就是按加入时间新到旧
// 		Find(&labors).
// 		Error; err != nil {
// 		return nil, err
// 	}

// 	return labors, nil
// }

// CreateLaborDivision 实现 LaborRepo 接口的 CreateLaborDivision 方法
func (r *laborRepoImpl) CreateLaborDivision(labor *dbmodels.ProjectLaborDivision) error {
	if labor == nil {
		return errors.New("labor 不能为 nil")
	}

	return r.handle.Model(&dbmodels.ProjectLaborDivision{}).
		Create(labor).Error
}

// TeamMemberRepo 接口定义了成员仓库的基本操作
type TeamMemberRepo interface {
	// SelectUserBasicPageWithParams 根据参数获取指定的成员，支持分页（ID 顺序）
	SelectUserBasicPageWithParams(offset, limit int, queryParams *dtos.MemberSearchParams) ([]*dbmodels.TeamMember, error)

	// SelectTeamBasicByUserId 获取指定成员 ID 下加入的所有汉化组信息
	SelectTeamBasicByUserId(userId dbmodels.PrimaryKey) ([]*dbmodels.Team, error)

	// SelectByUserId 获取指定用户 ID 的成员信息
	// 注意：这个函数会预加载 FkUser 的基础字段
	SelectByUserId(userId dbmodels.PrimaryKey) (*dbmodels.TeamMember, error)
	// // SelectByIdBacth 批量获取指定的 ID 的成员信息
	// // 这个函数同样会预加载 FkUser 的基础字段
	// SelectByIdBatch(ids []dbmodels.PrimaryKey) ([]dbmodels.TeamMember, error)
}

// teamMemberRepoImpl 是 MemberRepo 的实现
type teamMemberRepoImpl struct {
	handle *gorm.DB
}

// NewTeamMemberRepo 创建一个新的 MemberRepo 实例
func NewTeamMemberRepo(db *gorm.DB) TeamMemberRepo {
	return &teamMemberRepoImpl{
		handle: db,
	}
}

// SelectUserBasicPageWithParams 实现 MemberRepo 接口的 SelectUserBasicPageWithParams 方法
func (r *teamMemberRepoImpl) SelectUserBasicPageWithParams(
	offset, limit int,
	queryParams *dtos.MemberSearchParams,
) ([]*dbmodels.TeamMember, error) {
	var members []*dbmodels.TeamMember

	query := r.buildQueryWithParams(r.handle.Model(&dbmodels.TeamMember{}), queryParams)

	if err := query.
		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
			return db.Select(kUserBasicFields) // 只选择需要的字段
		}).
		Select("id, team_id, user_id, role, created_at, updated_at"). // 只选择需要的字段以提高性能
		Offset(offset).
		Limit(limit).
		Find(&members).
		Error; err != nil {
		return nil, err
	}

	return members, nil
}

// SelectByUserId 实现 MemberRepo 接口的 SelectByUserId 方法
func (r *teamMemberRepoImpl) SelectByUserId(userId dbmodels.PrimaryKey) (*dbmodels.TeamMember, error) {
	var member dbmodels.TeamMember

	if err := r.handle.Model(&dbmodels.TeamMember{}).
		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
			return db.Select(kUserBasicFields) // 只选择需要的字段
		}).
		Where("user_id = ?", userId).
		First(&member).
		Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// buildQueryWithParams 根据参数构建 WHERE 子句
func (r *teamMemberRepoImpl) buildQueryWithParams(base *gorm.DB, queryParams *dtos.MemberSearchParams) *gorm.DB {
	if queryParams == nil {
		return base
	}

	if queryParams.Nickname != nil {
		base = base.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", *queryParams.Nickname))
	}

	// 对于 QQ 号的 or 条件处理会复杂一些
	if queryParams.QqNumber != nil {
		switch queryParams.Nickname {
		case nil:
			// 如果之前没有添加 nickname 的 WHERE，则直接添加 WHERE
			base = base.Where("qq_number = ?", *queryParams.QqNumber)
		default:
			// 否则，使用 OR 连接
			base = base.Or("qq_number = ?", *queryParams.QqNumber)
		}
	}

	return base
}

// SelectTeamBasicByUserId 实现 MemberRepo 接口的 SelectTeamBasicByUserId 方法
func (r *teamMemberRepoImpl) SelectTeamBasicByUserId(userId dbmodels.PrimaryKey) ([]*dbmodels.Team, error) {
	var teams []*dbmodels.Team

	if err := r.handle.Model(&dbmodels.Team{}).
		Joins("RIGHT JOIN team_members ON teams.id = team_members.team_id").
		Where("team_members.user_id = ?", userId).
		Find(&teams).Error; err != nil {
		return nil, err
	}

	return teams, nil
}

// // SelectById 实现 MemberRepo 接口的 SelectById 方法
// func (r *teamMemberRepoImpl) SelectById(memberId dbmodels.PrimaryKey) (*dbmodels.TeamMember, error) {
// 	var member dbmodels.TeamMember

// 	if err := r.handle.Model(&dbmodels.TeamMember{}).
// 		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
// 			return db.Select(kUserBasicFields) // 只选择需要的字段
// 		}). // 全部预加载 FkUser 关联
// 		Where("id = ?", memberId).
// 		First(&member).
// 		Error; err != nil {
// 		return nil, err
// 	}

// 	return &member, nil
// }

// // SelectByIdBatch 实现 MemberRepo 接口的 SelectByIdBatch 方法
// // 接收一个 PrimaryKey 类型的 ID 切片
// func (r *teamMemberRepoImpl) SelectByIdBatch(ids []dbmodels.PrimaryKey) ([]dbmodels.TeamMember, error) {
// 	var members []dbmodels.TeamMember

// 	if len(ids) == 0 {
// 		// 如果 ID 列表为空，直接返回空切片
// 		return []dbmodels.TeamMember{}, nil
// 	}

// 	if err := r.handle.Model(&dbmodels.TeamMember{}).
// 		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
// 			return db.Select(kUserBasicFields) // 只选择需要的字段
// 		}).                      // 全部预加载 FkUser 关联
// 		Where("id IN (?)", ids). // 使用 IN 查询来批量查询
// 		Find(&members).
// 		Error; err != nil {
// 		return nil, err
// 	}

// 	return members, nil
// }
