package transformer

import (
	"poplargrid/internal/shared/dbmodels"
)

// Transformer 是数据转换器的接口
type Transformer struct {
}

// NewTransformer 构造一个新的 Transformer 实例
func NewTransformer() *Transformer {
	return &Transformer{}
}

// ProjectSetsToWorksets 将尨译的 project-set 格式转化成 Workset 格式
func (t *Transformer) ProjectSetsToWorksets(projsets []MoetranProjSet,
) ([]*dbmodels.Workset, error) {
	var worksets []*dbmodels.Workset
	for _, projset := range projsets {
		title := projset.Name
		if title == "default" || title == "" {
			title = "未分组"
		}

		worksets = append(worksets, &dbmodels.Workset{
			BaseModel: dbmodels.BaseModel{
				CreatedAt: projset.CreateTime.Time,
				UpdatedAt: projset.EditTime.Time,
			},
			Title:     projset.Name,
			MoetranId: projset.Id,
		})
	}

	return worksets, nil
}

// ProjsToWorks 从尨译的 project 信息提取出 Work 格式信息
func (t *Transformer) ProjsToWorks(projects []MoetranProj,
) ([]*dbmodels.Work, error) {
	// 将 MoetranProj 部分信息转化为 dbmodels.Work
	var works []*dbmodels.Work
	for _, project := range projects {
		works = append(works, &dbmodels.Work{
			BaseModel: dbmodels.BaseModel{
				CreatedAt: project.CreateTime.Time,
				UpdatedAt: project.EditTime.Time,
			},
			Title:       project.Name,
			MoetranId:   project.Id,
			Description: project.Intro,
		})
	}

	return works, nil
}

// // ProjsToProjects 将尨译的 project 信息提取出 Project 格式信息
// func (t *Transformer) ProjsToProjects(projects []crawler.MoetranProj,
// ) ([]*dbmodels.Project, error) {
// 	// 将 MoetranProj 部分信息转化为 dbmodels.Project
// 	var dbProjects []*dbmodels.Project
// 	for _, project := range projects {
// 		// 尝试获取 Legacy ID
// 		legacyId, _ := sExtractLegacyId(project.Name)
// 		// 提取纯净的标题
// 		cleanTitle := sExtractCleanTitle(project.Name)

// 		dbProjects = append(dbProjects, &dbmodels.Project{
// 			BaseModel: dbmodels.BaseModel{
// 				CreatedAt: project.CreateTime,
// 			},
// 			LegacyId: legacyId,
// 		})
// 	}

// 	return dbProjects, nil
// }
