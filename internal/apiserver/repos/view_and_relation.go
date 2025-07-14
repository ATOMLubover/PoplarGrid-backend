package repos

import (
	"poplargrid/internal/shared/dbmodels"

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

	// SelectProjectBasicPageIdDescByUserId 获取用户参与的项目列表，按 ID 倒序，支持分页
	SelectProjectBasicPageIdDescByUserId(userId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.Project, error)
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
		// 这里使用递归的 Preload 来加载成员和用户信息
		Preload("FkMember", func(dbMem *gorm.DB) *gorm.DB {
			return dbMem.
				Select("id", "user_id"). // Member 只需要 id、user_id
				Preload("FkMember.FkUser", func(dbUser *gorm.DB) *gorm.DB {
					return dbUser.Select("id, nickname") // User 只需要 id、nickname
				})
		}).
		Where("project_id = ?", projectId).
		Find(&labors).
		Error; err != nil {
		return nil, err
	}

	return labors, nil
}

// SelectProjectBasicPageIdDescByUserId 实现 LaborRepo 接口的 SelectProjectBasicPageIdDescByUserId 方法
func (r *laborRepoImpl) SelectProjectBasicPageIdDescByUserId(userId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.Project, error) {
	var projects []*dbmodels.Project

	// 为每个字段前加上表名限定
	var selectFields []string

	for _, field := range kProjectBasicFields {
		selectFields = append(selectFields, "projects."+field)
	}

	if err := r.handle.Model(&dbmodels.Project{}).
		Select(selectFields).
		// 使用 INNER JOIN 连接 project_labor_divisions 表，选出指定用户参与的项目
		// 如果用户没有参与任何项目，则不会返回任何结果
		Joins("JOIN project_labor_divisions ON project_labor_divisions.project_id = projects.id").
		Where("project_labor_divisions.user_id = ?", userId).
		Order("projects.id DESC").
		Offset(offset).
		Limit(limit).
		Find(&projects).
		Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// TeamMemberRepo 接口定义了成员仓库的基本操作
type TeamMemberRepo interface {
	// SelectUserBasicPage 获取指定团队 ID 的成员列表，支持分页（ID 顺序）
	SelectUserBasicPage(teamId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.TeamMember, error)

	// SelectTeamBasicByUserId 获取指定成员 ID 下加入的所有汉化组信息
	SelectTeamBasicByUserId(userId dbmodels.PrimaryKey) ([]*dbmodels.Team, error)

	// // SelectById 获取指定成员 ID 的成员信息
	// // 注意：这个函数会预加载 FkUser 的基础字段
	// SelectById(id dbmodels.PrimaryKey) (*dbmodels.TeamMember, error)
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

// SelectUserBasicPage 实现 MemberRepo 接口的 SelectUserBasicPage 方法
func (r *teamMemberRepoImpl) SelectUserBasicPage(teamId dbmodels.PrimaryKey, offset, limit int) ([]*dbmodels.TeamMember, error) {
	var members []*dbmodels.TeamMember

	if err := r.handle.Model(&dbmodels.TeamMember{}).
		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
			return db.Select(kUserBasicFields) // 只选择需要的字段
		}).
		Select("id, team_id, user_id, role, created_at, updated_at"). // 只选择需要的字段以提高性能
		Where("team_id = ?", teamId).
		Offset(offset).
		Limit(limit).
		Find(&members).
		Error; err != nil {
		return nil, err
	}

	return members, nil
}

// SelectTeamBasicByUserId 实现 MemberRepo 接口的 SelectTeamBasicByUserId 方法
func (r *teamMemberRepoImpl) SelectTeamBasicByUserId(userId dbmodels.PrimaryKey) ([]*dbmodels.Team, error) {
	var teams []*dbmodels.Team

	if err := r.handle.Model(&dbmodels.TeamMember{}).
		Select("Team.*").
		// 使用 LEFT JOIN 连接 teams 表，保证即使用户没有加入任何团队也能返回 NULL 结果
		Joins("LEFT JOIN teams ON teams.id = team_members.team_id").
		Where("team_members.user_id = ?", userId).
		Find(&teams).Error; err != nil {
		return nil, err
	}

	return teams, nil
}

// SelectById 实现 MemberRepo 接口的 SelectById 方法
func (r *teamMemberRepoImpl) SelectById(memberId dbmodels.PrimaryKey) (*dbmodels.TeamMember, error) {
	var member dbmodels.TeamMember

	if err := r.handle.Model(&dbmodels.TeamMember{}).
		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
			return db.Select(kUserBasicFields) // 只选择需要的字段
		}). // 全部预加载 FkUser 关联
		Where("id = ?", memberId).
		First(&member).
		Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// SelectByIdBatch 实现 MemberRepo 接口的 SelectByIdBatch 方法
// 接收一个 PrimaryKey 类型的 ID 切片
func (r *teamMemberRepoImpl) SelectByIdBatch(ids []dbmodels.PrimaryKey) ([]dbmodels.TeamMember, error) {
	var members []dbmodels.TeamMember

	if len(ids) == 0 {
		// 如果 ID 列表为空，直接返回空切片
		return []dbmodels.TeamMember{}, nil
	}

	if err := r.handle.Model(&dbmodels.TeamMember{}).
		Preload("FkUser", func(db *gorm.DB) *gorm.DB {
			return db.Select(kUserBasicFields) // 只选择需要的字段
		}).                      // 全部预加载 FkUser 关联
		Where("id IN (?)", ids). // 使用 IN 查询来批量查询
		Find(&members).
		Error; err != nil {
		return nil, err
	}

	return members, nil
}
