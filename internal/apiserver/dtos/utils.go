package dtos

// ProjectOverallStatus 是项目所有状态的位掩码类型别名
type ProjectOverallStatus uint

// 项目阶段状态枚举
const (
	PROJECT_STATUS_UNSET       uint = 0b00 // 未开始
	PROJECT_STATUS_IN_PROGRESS uint = 0b01 // 进行中
	PROJECT_STATUS_COMPLETED   uint = 0b10 // 已完成
	PROJECT_STATUS_NOT_USED    uint = 0b11 // 指定不参与搜索
)

// 各阶段状态的掩码位移常量
const (
	PROJECT_STATUS_TRANSLATING_SHIFT  uint = iota * 2 // 翻译状态的位移 - 位移 0 (占用 0, 1 位)
	PROJECT_STATUS_PROOFREADING_SHIFT                 // 校对状态的位移 - 位移 2 (占用 2, 3 位)
	PROJECT_STATUS_LETTERING_SHIFT                    // 字幕制作状态的位移 - 位移 4 (占用 4, 5 位)
	PROJECT_STATUS_REVIEWING_SHIFT                    // 审核状态的位移 - 位移 6 (占用 6, 7 位)
)

// 各阶段状态对应的掩码常量 (用于写入和读取)
const (
	PROJECT_STATUS_TRANSLATE_MASK ProjectOverallStatus = 0b000000011 // 二进制 3
	PROJECT_STATUS_PROOF_MASK     ProjectOverallStatus = 0b000001100 // 二进制 12
	PROJECT_STATUS_LETTER_MASK    ProjectOverallStatus = 0b000110000 // 二进制 48
	PROJECT_STATUS_REVIEW_MASK    ProjectOverallStatus = 0b011000000 // 二进制 192
)

// SetTranslatingStatus 更新 ProjectOverallStatus 中的翻译状态位
func (ps *ProjectOverallStatus) SetTranslatingStatus(status uint) {
	*ps &= ^PROJECT_STATUS_TRANSLATE_MASK
	*ps |= (ProjectOverallStatus(status) << PROJECT_STATUS_TRANSLATING_SHIFT)
}

// GetTranslatingStatus 从 ProjectOverallStatus 中获取翻译状态
func (ps ProjectOverallStatus) GetTranslatingStatus() uint {
	return uint((ps & PROJECT_STATUS_TRANSLATE_MASK) >> PROJECT_STATUS_TRANSLATING_SHIFT)
}

// SetProofreadingStatus 更新 ProjectOverallStatus 中的校对状态位
func (ps *ProjectOverallStatus) SetProofreadingStatus(status uint) {
	*ps &= ^PROJECT_STATUS_PROOF_MASK
	*ps |= (ProjectOverallStatus(status) << PROJECT_STATUS_PROOFREADING_SHIFT)
}

// GetProofreadingStatus 从 ProjectOverallStatus 中获取校对状态
func (ps ProjectOverallStatus) GetProofreadingStatus() uint {
	return uint((ps & PROJECT_STATUS_PROOF_MASK) >> PROJECT_STATUS_PROOFREADING_SHIFT)
}

// SetLetteringStatus 更新 ProjectOverallStatus 中的字幕制作状态位
func (ps *ProjectOverallStatus) SetLetteringStatus(status uint) {
	*ps &= ^PROJECT_STATUS_LETTER_MASK
	*ps |= (ProjectOverallStatus(status) << PROJECT_STATUS_LETTERING_SHIFT)
}

// GetLetteringStatus 从 ProjectOverallStatus 中获取字幕制作状态
func (ps ProjectOverallStatus) GetLetteringStatus() uint {
	return uint((ps & PROJECT_STATUS_LETTER_MASK) >> PROJECT_STATUS_LETTERING_SHIFT)
}

// SetReviewingStatus 更新 ProjectOverallStatus 中的审核状态位
func (ps *ProjectOverallStatus) SetReviewingStatus(status uint) {
	*ps &= ^PROJECT_STATUS_REVIEW_MASK
	*ps |= (ProjectOverallStatus(status) << PROJECT_STATUS_REVIEWING_SHIFT)
}

// GetReviewingStatus 从 ProjectOverallStatus 中获取审核状态
func (ps ProjectOverallStatus) GetReviewingStatus() uint {
	return uint((ps & PROJECT_STATUS_REVIEW_MASK) >> PROJECT_STATUS_REVIEWING_SHIFT)
}

// 获取列表时 sort 参数的可选值
const (
	SORT_ID_DESC         = iota // 按照 ID 倒序
	SORT_UPDATED_AT_DESC        // 按照更新时间倒序
)
