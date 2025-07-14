package services

import (
	"poplargrid/internal/api_server/repository"
	"poplargrid/internal/shared/dbmodels"
)

// ProjectService 提供项目相关的服务
type ProjectService struct {
	projRepo *repository.ProjectsRepo
}

// NewProjectService 创建一个新的 ProjectService 实例
func NewProjectService(
	projRepo *repository.ProjectsRepo,
) *ProjectService {
	return &ProjectService{
		projRepo: projRepo,
	}
}

// GetProjectById 根据项目 ID 获取项目的完整信息
// 如果不存在则返回 nil 和错误
func (s *ProjectService) GetProjectById(projectId uint) (*dbmodels.Project, error) {
	return s.projRepo.SelectFullById(projectId)
}

// GetAllProjects 获取所有项目的完整信息
// 返回所有项目的完整信息
func (s *ProjectService) GetAllProjects(offset, num int) (
	[]*dbmodels.Project,
	map[dbmodels.PrimaryKey][]*dbmodels.Tag,
	map[dbmodels.PrimaryKey][]*dbmodels.User,
	error) {
	return s.projRepo.SelectAllToSlice(offset, num)
}
