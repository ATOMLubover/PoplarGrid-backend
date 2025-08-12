package crawler

import (
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/shared/models"
)

// teamMoetranToPoplar 将尨译的汉化组信息转化为 PoplarGrid 的汉化组模型
func teamMoetranToPoplar(moetranTeams []apiclient.TeamDTO) []models.Team {
	poplarTeams := make([]models.Team, len(moetranTeams))

	for i, team := range moetranTeams {
		poplarTeams[i] = models.Team{
			MoetranId: team.ID,
			Name:      team.Name,
		}
	}

	return poplarTeams
}

// setMoetranToPoplar 将尨译的项目集信息转化为 PoplarGrid 的作品集模型
func setMoetranToPoplar(team *models.Team, moetranSets []apiclient.ProjectSetDTO) []models.Workset {
	poplarSets := make([]models.Workset, len(moetranSets))

	for i, set := range moetranSets {
		poplarSets[i] = models.Workset{
			MoetranId: set.ID,
			TeamId:    team.Id,
			Name:      set.Name,
		}
	}

	return poplarSets
}

// projectMoetranToPoplar 将尨译的项目信息转化为 PoplarGrid 的项目模型
func projectMoetranToPoplar(memberID uint, workset *models.Workset, moetranProjects []apiclient.ProjectDTO) []models.Project {
	poplarProjects := make([]models.Project, len(moetranProjects))

	for i, project := range moetranProjects {
		poplarProjects[i] = models.Project{
			MoetranId:   project.ID,
			WorksetId:   workset.Id,
			Title:       project.Name,
			Description: project.Intro,
			PrincipalId: models.PKey(memberID),
		}
	}

	return poplarProjects
}
