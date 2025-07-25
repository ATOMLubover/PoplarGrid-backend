package models

import "gorm.io/gorm"

// WorksetStat 映射到数据库的 project_stats_mv 物化视图
// 注意：这个 struct 只用于读取，不用于 GORM 的自动迁移或写入操作
type WorksetStat struct {
	// WorksetId 现在是每行的唯一标识，用于并发刷新
	WorksetId PKey `gorm:"column:workset_id;primaryKey"`

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
	Published int `gorm:"column:published"`
}

// TableName 返回 WorkStatsMv 的物化视图名
func (*WorksetStat) TableName() string {
	return "workset_stats_mv"
}

// 外部调用的指针，避免每一次都创建
var emptyWorksetStat = &WorksetStat{}

// GetWorkStatsMv 获取一个空的 WorkStat 实例
func GetWorkStatsMv() *WorksetStat {
	return emptyWorksetStat
}

// Select 查询 WorkStatsMv 的单行记录
func (*WorksetStat) Select(hdl *gorm.DB, id PKey) (*WorksetStat, error) {
	var stats WorksetStat

	if err := hdl.
		Model(&WorksetStat{}).
		Where("workset_id = ?", id).
		First(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
