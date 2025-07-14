package dbmodels

// ProjectStats 映射到数据库的 project_stats_mv 物化视图
// 注意：这个 struct 只用于读取，不用于 GORM 的自动迁移或写入操作
type ProjectStats struct {
	// WorksetId 现在是每行的唯一标识
	WorksetId PrimaryKey `gorm:"column:workset_id;primaryKey"`

	// 总项目数
	Total int `gorm:"column:total"`

	// 翻译统计
	NotTranslating int `gorm:"column:not_translating"`
	Translating    int `gorm:"column:translate_in_progress"`
	Translated     int `gorm:"column:translate_completed"`

	// 校对统计
	NotProoving int `gorm:"column:not_prooving"`
	Prooving    int `gorm:"column:proof_in_progress"`
	Prooved     int `gorm:"column:proof_completed"`

	// 排版统计
	NotLettering int `gorm:"column:not_lettering"`
	Lettering    int `gorm:"column:letter_in_progress"`
	Letterred    int `gorm:"column:letter_completed"`

	// 审核统计
	NotReviewing int `gorm:"column:not_reviewing"`
	Reviewing    int `gorm:"column:review_in_progress"`
	Reviewed     int `gorm:"column:review_completed"`

	// 发布状态
	NotPublished int `gorm:"column:not_published"`
	Published    int `gorm:"column:published"`
}

// 注意：这里是物化视图名
// 刷新方法；REFRESH MATERIALIZED VIEW CONCURRENTLY project_stats_mv;
func (ProjectStats) TableName() string {
	return "project_stats_mv"
}
