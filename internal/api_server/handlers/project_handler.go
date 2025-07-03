package handlers

import (
	"poplargrid/internal/api_server/dtos"
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
)

// ProjectHandler 处理项目相关的请求
type ProjectHandler struct {
	ProjService *services.ProjectService
}

// ProjListPage godoc
// @Summary      获取项目列表分页
// @Description  注意当列表为空返回值是 null 而非 []
// @Param        pageSerial query int false "页码，默认值为 1"
// @Param        pageSize query int false "每页数量，默认值为 10"
// @Tags         project
// @Produce      json
// @Success      200 {object} []dtos.ProjectFullInfo
// @Router       /proj/list [get]
func (h *ProjectHandler) ProjListPage(ctx iris.Context) {
	// 获取分页参数
	pageSerial := ctx.URLParamIntDefault("pageSerial", 1)
	pageSize := ctx.URLParamIntDefault("pageSize", 10)

	// 确保分页参数合法
	if pageSerial < 1 {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"error":      "pageSerial非法",
			"pageSerial": pageSerial,
		})
		return
	}
	if pageSize < 1 {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"error":      "pageSize非法",
			"pageSerial": pageSerial,
		})
		return
	}

	// 处理分页参数为 offset 和 limit
	offset := (pageSerial - 1) * pageSize
	limit := pageSize

	// 调用服务获取项目列表
	projs, tagMap, laborMap, err :=
		h.ProjService.GetAllProjects(offset, limit)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error":   "获取项目列表失败",
			"details": err.Error(),
		})
		return
	}

	// 成功返回项目详情列表
	var responseSlice []*dtos.ProjectFullInfo

	for _, proj := range projs {
		// 加载项目对应 work 的 tag
		tags := make([]string, 0)
		if tagList, ok := tagMap[proj.WorkId]; ok {
			for _, tag := range tagList {
				tags = append(tags, tag.Name)
			}
		}

		// 加载项目分工成员
		laborDivision := make([]dtos.ProjectMemberLabor, 0)
		if laborList, ok := laborMap[proj.Id]; ok {
			for _, member := range laborList {
				laborDivision = append(laborDivision, dtos.ProjectMemberLabor{
					MemberName: member.Nickname,
					LaborRole:  member.Labors, // 这里是借用了 Member 的 Labors 字段传递分工角色
				})
			}
		}

		// 构建项目完整信息 DTO
		responseSlice = append(responseSlice, &dtos.ProjectFullInfo{
			ProjectId: uint(proj.Id),
			Title:     proj.FkWork.Title,
			Team: dtos.TeamFullInfo{
				TeamId:    uint(proj.TeamId),
				MoetranId: proj.FkTeam.MoetranId,
				TeamName:  proj.FkTeam.Name,
			},
			Workset: dtos.WorksetFullInfo{
				WorksetId: uint(proj.WorksetId),
				MoetranId: proj.FkWorkset.MoetranId,
				Title:     proj.FkWorkset.Title,
			},
			Work: dtos.WorkFullInfo{
				WorkId:      uint(proj.WorkId),
				MoetranId:   proj.FkWork.MoetranId,
				Title:       proj.FkWork.Title,
				Description: proj.FkWork.Description,
				Tags:        tags,
			},
			Status:        proj.Status,
			Urgency:       proj.Urgency,
			LegacyId:      proj.LegacyId,
			LaborDivision: laborDivision,
		})
	}

	ctx.JSON(responseSlice)
}
