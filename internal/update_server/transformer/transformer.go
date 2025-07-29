package transformer

import (
	"poplargrid/internal/shared/models"
	"poplargrid/internal/update_server/apidto"
	"regexp"
	"strconv"
)

// UsersToLocalUsers 将尨译的用户信息转化为本地 Member 模型
func UsersToLocalUsers(team *models.Team, users []apidto.MoetranUser) ([]*models.User, error) {
	var localusers []*models.User

	for _, user := range users {
		localusers = append(localusers, &models.User{
			Nickname:  user.Name,
			MoetranId: user.Id,

			Email: user.Id, // 这里使用 user.Id 作为 Email 占位，防止新创建时违反 unique
		})
	}

	return localusers, nil
}

// LocalUsersToTeamMembers 将尨译的用户信息转化为本地 TeamMember 模型
func LocalUsersToTeamMembers(team *models.Team, users []*models.User) ([]*models.Member, error) {
	var members []*models.Member

	for _, user := range users {
		members = append(members, &models.Member{
			UserId: user.Id, // 使用转换后的用户 ID
			TeamId: team.Id, // 使用传入的 team 的 ID
		})
	}

	return members, nil
}

// ProjSetsToWorksets 将尨译的 project-set 格式转化成 Workset 格式
func ProjSetsToWorksets(team *models.Team, projsets []apidto.MoetranProjSet) (
	[]*models.Workset, error) {
	var worksets []*models.Workset

	for _, projset := range projsets {
		name := projset.Name
		if name == "default" || name == "" {
			name = "未分组"
		}

		worksets = append(worksets, &models.Workset{
			BaseModel: models.BaseModel{
				CreatedAt: projset.CreateTime.Time,
				UpdatedAt: projset.EditTime.Time,
			},
			Name:      name,
			MoetranId: projset.Id,
			TeamId:    team.Id, // 使用传入的 team 的 ID
		})
	}

	return worksets, nil
}

// ProjsToWorks 从尨译的 project 信息提取出 Project 格式信息
// workset 是辅助处理的作品集信息，为当前 projs 所在的作品集
func ProjsToProjects(projs []apidto.MoetranProj, workset *models.Workset) (
	[]*models.Project, error) {
	// 将 MoetranProj 信息转化为 dbmodel.Project
	var projects []*models.Project

	for _, proj := range projs {
		projects = append(projects, &models.Project{
			BaseModel: models.BaseModel{
				CreatedAt: proj.CreateTime.Time,
				UpdatedAt: proj.EditTime.Time,
			},
			Title:     sExtractCleanTitle(proj.Name),
			MoetranId: proj.Id,
			LegacyId:  sExtractLegacyId(proj.Name),
			WorksetId: workset.Id, // 使用传入的 workset 的 ID

			PrincipalId: 1, // 默认负责人 ID 为 1
		})
	}

	return projects, nil
}

// ============== 辅助函数 ==============

// sExtractLegacyId 从标题字符串中提取出 Legacy ID 部分
// 返回值为提取出的 Legacy ID 和一个布尔值，表示是否成功提取
func sExtractLegacyId(title string) uint {
	// 使用正则表达式匹配标题中的 Legacy ID
	// Legacy ID 的格式为【数字】
	re := regexp.MustCompile(`【(\d+)】`)

	matches := re.FindStringSubmatch(title)

	if len(matches) <= 0 {
		// 如果没有匹配到任何内容，返回 0
		return 0
	}

	if len(matches) > 1 {
		// 获取匹配的第一个子字符串，即【】之间的数字部分
		legacyIdStr := matches[1]

		// 将字符串数字转换为整数
		legacyId, err := strconv.ParseUint(
			legacyIdStr, 10, 32)
		if err != nil {
			return 0
		}

		return uint(legacyId)
	}

	// 如果只匹配到 Legacy ID 的整体格式，但没有数字部分，返回 0
	// 虽然这种情况不太可能，因为正则表达式肯定会匹配到数字
	return 0
}

// sExtractFullName 去除标题字符串中的 Legacy ID 部分
// 返回值为纯净的标题字符串
func sExtractCleanTitle(title string) string {
	// 使用正则表达式匹配标题中的 Legacy ID
	// Legacy ID 整体的格式为【一个数字】
	re := regexp.MustCompile(`【\d+】`)

	// 替换匹配到的 Legacy ID 部分为空字符串
	return re.ReplaceAllString(title, "")
}
