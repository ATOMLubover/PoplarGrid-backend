package dtos

// CreateProjectInfo 定义了创建项目所需的信息
type CreateProjectInfo struct {
	Title       string // 项目标题
	Description string // 项目简介
	WorksetId   uint   // 作品集 ID

	CreatorUserId uint // 创建者的用户 ID
	AllowAutoJoin bool // 允许加入的权限
	IsHidden      bool // 是否隐藏项目
}
