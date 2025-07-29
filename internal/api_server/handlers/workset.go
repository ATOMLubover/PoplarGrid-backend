package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// WorksetInfo 定义了工作集的基本信息 DTO
type WorksetInfo struct {
	// 工作集 ID
	Id uint `json:"id"`
	// 工作集名称
	Name string `json:"name"`
	// 工作集描述
	Description string `json:"description,omitempty"`
	// 工作集对应的龙译 ID
	MoetranId string `json:"moetran_id,omitempty"`
	// 工作集的创建时间
	CreatedAt string `json:"created_at"`
}

// CreateWorksetRequest 定义了创建工作集的请求参数
type CreateWorksetRequest struct {
	// 工作集名称
	Name string `json:"name" validate:"required"`
	// 工作集描述
	Description string `json:"description"`
	// 所属团队 ID
	TeamId uint `json:"team_id" validate:"required"`
}

// WorksetStats 定义了作品集项目的统计信息 DTO
type WorksetStats struct {
	// 作品集 ID
	WorksetId uint `json:"workset_id"`
	// 总项目数量
	TotalCount int `json:"total_project_count"`

	// 未开始翻译的项目数量
	NotTranslatingCount int `json:"not_translating_count"`
	// 翻译中的项目数量
	OnTranslatingCount int `json:"on_translating_count"`
	// 翻译完成的项目数量
	TranslatedCount int `json:"translated_count"`

	// 未开始校对的项目数量
	NotProofreadingCount int `json:"not_proofreading_count"`
	// 校对中的项目数量
	OnProofreadingCount int `json:"on_proofreading_count"`
	// 校对完成的项目数量
	ProofreadCount int `json:"proofread_count"`

	// 未开始嵌字的项目数量
	NotLetteringCount int `json:"not_lettering_count"`
	// 嵌字中的项目数量
	OnLetteringCount int `json:"on_lettering_count"`
	// 嵌字完成的项目数量
	LetteredCount int `json:"lettered_count"`

	// 未开始审核的项目数量
	NotReviewingCount int `json:"not_reviewing_count"`
	// 审核中的项目数量
	OnReviewingCount int `json:"on_reviewing_count"`
	// 审核完成的项目数量
	ReviewedCount int `json:"reviewed_count"`

	// 已发布的项目数量
	PublishedCount int `json:"published_count"`
}

// RouteWorksetHandler 注册工作集相关的路由
func RouteWorksetHandler(root *mvc.Application) {
	// 注册路由组
	root.Party("/worksets").
		Handle(new(WorksetHandler))
}

// WorksetHandler 处理工作集相关的请求
type WorksetHandler struct {
	WorksetService services.WorksetService
}

// BeforeActivation 在控制器激活前注册路由
func (w *WorksetHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "List")
	b.Handle("GET", "/{id:uint}/stats", "Stats")
	b.Handle("POST", "/", "Create")
}

// List godoc
//
// @Summary 	获取工作集列表
// @Description 注意当列表为空，会返回 null 而不是空数组
//
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		user_id query string true "汉化组 ID"
//
// @Tags 		workset
// @Produce 	json
// @Success	 	200 {object} []WorksetInfo
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/worksets [get]
func (h *WorksetHandler) List(ctx iris.Context) {
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	// 调用服务获取工作集列表
	worksets, err := h.WorksetService.GetWorksets(&services.WorksetListParams{
		Offset: (pageSerial - 1) * pageSize,
		Limit:  pageSize,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error: "获取工作集列表失败",
		})
		return
	}

	// 将 service 的 WorksetInfo 转换为 DTO
	worksetInfos := make([]*WorksetInfo, 0, len(worksets))

	for _, ws := range worksets {
		worksetInfos = append(worksetInfos, &WorksetInfo{
			Id:          ws.Id,
			Name:        ws.Name,
			Description: ws.Description,
			MoetranId:   ws.MoetranId,
		})
	}

	ctx.JSON(worksetInfos)
}

// Stats godoc
//
// @Summary 	获取特定作品集项目统计信息
// @Description 获取所有项目的统计信息，包括总数、翻译、校对、嵌字，审核、发布对应数量等
//
// @Param 		id path int true "作品集 ID"
//
// @Tags 		workset
// @Produce 	json
// @Success	 	200 {object} WorksetStats
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/api/worksets/{id}/stats [get]
func (h *WorksetHandler) Stats(ctx iris.Context) {
	// 获取 workset_id 参数
	worksetId, err := ctx.Params().GetUint("id")
	if err != nil || worksetId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "workset_id 必须是一个明确给出的正整数",
		})
		return
	}

	// 调用服务层获取统计数据
	stats, err := h.WorksetService.GetWorksetStats(worksetId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ctx.JSON(ErrorResponse{
			Error:  "获取特定作品集的项目统计信息失败",
			Detail: err.Error(),
		}))
		return
	}

	// 将 stats 转换成 DTO
	ctx.JSON(WorksetStats{
		WorksetId:            stats.WorksetId,
		TotalCount:           stats.TotalCount,
		NotTranslatingCount:  stats.NotTranslatingCount,
		OnTranslatingCount:   stats.TranslatedCount,
		TranslatedCount:      stats.TranslatedCount,
		NotProofreadingCount: stats.NotProovingCount,
		OnProofreadingCount:  stats.ProovingCount,
		ProofreadCount:       stats.ProovedCount,
		NotLetteringCount:    stats.NotLetteringCount,
		OnLetteringCount:     stats.LetteringCount,
		LetteredCount:        stats.LetteredCount,
		NotReviewingCount:    stats.NotReviewingCount,
		OnReviewingCount:     stats.ReviewingCount,
		ReviewedCount:        stats.ReviewedCount,
		PublishedCount:       stats.PublishedCount,
	})
}

// Create godoc
//
// @Summary 	创建新的工作集
// @Description 创建一个新的工作集，必须提供名称和所属团队 ID
//
// @Accept 		json
// @Param 		body_params body CreateWorksetRequest true "创建工作集请求参数"
//
// @Tags 		workset
// @Produce 	json
// @Success 	200 {object} SuccessResponse "创建成功"
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/worksets [post]
func (h *WorksetHandler) Create(ctx iris.Context) {
	// 读取请求体中的工作集信息
	var req CreateWorksetRequest

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的请求参数",
		})
		return
	}

	// 调用服务层创建工作集
	info, err := h.WorksetService.CreateWorkset(&services.CreateWorksetParams{
		Name:        req.Name,
		Description: req.Description,
		TeamId:      req.TeamId,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "创建工作集失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(SuccessResponse{
		Message: info.Message,
		Detail: &iris.Map{
			"workset_id": info.WorksetId,
			"moetran_id": info.MoetranId,
		},
	})
}
