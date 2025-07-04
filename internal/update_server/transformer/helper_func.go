package transformer

import (
	"regexp"
	"strconv"
)

// ExtractLegacyId 从标题字符串中提取出 Legacy ID 部分
// 返回值为提取出的 Legacy ID 和一个布尔值，表示是否成功提取
func sExtractLegacyId(title string) (uint, bool) {
	// 使用正则表达式匹配标题中的 Legacy ID
	// Legacy ID 的格式为【数字】
	re := regexp.MustCompile(`【(\d+)】`)

	matches := re.FindStringSubmatch(title)

	if len(matches) <= 0 {
		// 如果没有匹配到任何内容，返回 false
		return 0, false
	}

	if len(matches) > 1 {
		// 获取匹配的第一个子字符串，即【】之间的数字部分
		legacyIdStr := matches[1]

		// 将字符串数字转换为整数
		legacyId, err := strconv.ParseUint(
			legacyIdStr, 10, 32)
		if err != nil {
			return 0, false
		}

		return uint(legacyId), true
	}

	return 0, false
}

// ExtractFullName 去除标题字符串中的 Legacy ID 部分
// 返回值为纯净的标题字符串
func sExtractCleanTitle(title string) string {
	// 使用正则表达式匹配标题中的 Legacy ID
	// Legacy ID 整体的格式为【一个数字】
	re := regexp.MustCompile(`【\d+】`)

	// 替换匹配到的 Legacy ID 部分为空字符串
	return re.ReplaceAllString(title, "")
}
