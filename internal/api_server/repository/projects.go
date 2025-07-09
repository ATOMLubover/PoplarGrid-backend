package repository

import (
	"poplargrid/internal/shared/dbmodel"

	"gorm.io/gorm"
)

// ProjectsRepo 提供项目相关的数据库操作
type ProjectsRepo struct {
	DbCtx          *gorm.DB
	RelationTables *RelationTables // 关联表
}

// NewProjectRepo 构造一个新的 ProjectRepo 实例
func NewProjectsRepo(db *gorm.DB) *ProjectsRepo {
	return &ProjectsRepo{
		DbCtx:          db,
		RelationTables: NewRelationTables(db),
	}
}

// GetTable 获取项目表的上下文引用
func (r *ProjectsRepo) GetTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodel.Project{})
}

// SelectFullModelById 根据项目 ID 获取项目的完整信息
// 如果不存在则返回 nil 和错误
func (r *ProjectsRepo) SelectFullById(projectId uint) (*dbmodel.Project, error) {
	var project dbmodel.Project

	if err := r.GetTable().
		Where("id = ?", projectId).
		First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

// SelectAllToSlice 获取所有项目的完整信息
// 返回所有项目的完整信息
// 以及对应的作品标签和分工成员信息（以 projects 的 ID 为 key 的 map）
// 需要注意，这里 Member 的 role 是特属于这个项目的分工角色
func (r *ProjectsRepo) SelectAllToSlice(offset, num int) (
	[]*dbmodel.Project,
	map[dbmodel.PrimaryKey][]*dbmodel.Tag,
	map[dbmodel.PrimaryKey][]*dbmodel.Member,
	error,
) {
	var projects []*dbmodel.Project

	if err := r.GetTable().
		// 预查询 FkWorkset，并指定 Workset 需要的字段
		Preload("FkWorkset", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "moetran_id")
		}).
		Offset(offset).
		Limit(num).
		Find(&projects).Error; err != nil {
		return nil, nil, nil, err
	}

	// 如果 projects 为空，返回空切片而非 nil
	// 这样可以避免调用方需要处理 nil 的情况
	if len(projects) == 0 {
		return []*dbmodel.Project{},
			map[dbmodel.PrimaryKey][]*dbmodel.Tag{},
			map[dbmodel.PrimaryKey][]*dbmodel.Member{},
			nil
	}

	// 然后开始从关联表内中收集项目 tag
	projIds := make([]dbmodel.PrimaryKey, 0)
	// 收集所有作品 ID
	for _, project := range projects {
		projIds = append(projIds, project.Id)
	}
	// 然后获取所有作品对应的 tag
	projTagMap, err := r.sMapProjectTags(projIds)
	if err != nil {
		return nil, nil, nil, err
	}

	// 收集所有项目的分工
	projectIds := make([]dbmodel.PrimaryKey, 0)
	// 收集所有项目 ID
	for _, project := range projects {
		projectIds = append(projectIds, dbmodel.PrimaryKey(project.Id))
	}
	// 然后获取所有项目对应的分工成员
	laborMap, err := r.sMapLaborDiv(projectIds)
	if err != nil {
		return nil, nil, nil, err
	}

	// 最后返回获取的值
	return projects, projTagMap, laborMap, nil
}

// ================= 辅助函数 =================

// sMapWorkTags 查询各个 project 对应的所有 tag
// 返回一个 map，key 为 projectId，value 为该 work 的所有 tag
func (r *ProjectsRepo) sMapProjectTags(
	projIds []dbmodel.PrimaryKey,
) (map[dbmodel.PrimaryKey][]*dbmodel.Tag, error) {
	if len(projIds) == 0 {
		// 如果没有提供任何 workId，直接返回空 map
		return map[dbmodel.PrimaryKey][]*dbmodel.Tag{}, nil
	}

	projTags := make([]*dbmodel.ProjectTag, 0)

	if err := r.RelationTables.GetProjectTagTable().
		Preload("FkTag", func(tx *gorm.DB) {
			tx.Select("id", "name", "description") // 仅查询 tag 必要字段
		}). // 预加载关联的 tag
		Where("project_id IN (?)", projIds).
		Find(&projTags).Error; err != nil {
		return nil, err
	}

	// 向 map 中填充作品 tag 信息
	// key 为 workId，value 为该作品的所有 tag
	tagMap := make(map[dbmodel.PrimaryKey][]*dbmodel.Tag)
	for _, projTag := range projTags {
		if projTag.FkTag.Id != 0 {
			// 确保 tag ID 有效
			tagMap[projTag.ProjectId] = append(tagMap[projTag.ProjectId], &projTag.FkTag)
		}
	}

	return tagMap, nil
}

// sMapLaborDiv 查询所有项目的分工信息
// 返回一个 map，key 为 projectId，value 为该项目的所有分工成员
func (r *ProjectsRepo) sMapLaborDiv(
	projectIDs []dbmodel.PrimaryKey,
) (map[dbmodel.PrimaryKey][]*dbmodel.Member, error) {
	if len(projectIDs) == 0 {
		// 如果没有提供任何 projectId，直接返回空 map
		return map[dbmodel.PrimaryKey][]*dbmodel.Member{}, nil
	}

	var plds []dbmodel.ProjectLaborDivision

	if err := r.RelationTables.GetProjectLaborDivisionTable().
		Preload("FkMember", func(tx *gorm.DB) {
			tx.Select("id", "nickname") // 仅查询 Member 必要字段
		}). // 预加载关联的 Member
		Where("project_id IN (?)", projectIDs).
		Find(&plds).Error; err != nil {
		return nil, err
	}

	// 向 map 中填充分工信息
	// key 为 projectId，value 为该项目的所有分工成员
	laborMap := make(map[dbmodel.PrimaryKey][]*dbmodel.Member)
	for _, pld := range plds {
		if pld.FkMember.Id != 0 {
			// 确保 Member 被成功加载
			pld.FkMember.Labors = pld.LaborRole
			laborMap[pld.ProjectId] = append(laborMap[pld.ProjectId], &pld.FkMember)
		}
	}

	return laborMap, nil
}
