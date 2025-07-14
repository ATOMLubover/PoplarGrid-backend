package repository

import (
	"fmt"
	"log/slog"
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 作品集 repo
type WorksetsRepo struct {
	// DbCtx 是数据库上下文
	DbCtx *gorm.DB
	// 日志器
	logger *slog.Logger // 如果需要日志记录，可以添加日志器
}

// NewWorksetsRepo 构造一个新的 WorksetsRepo 实例
func NewWorksetsRepo(db *gorm.DB) *WorksetsRepo {
	return &WorksetsRepo{
		DbCtx: db,
	}
}

// GetTable 获取 worksets 表的上下文引用
func (r *WorksetsRepo) GetTable() *gorm.DB {
	return r.DbCtx.Model(&dbmodels.Workset{})
}

// BulkUpsert 批量更新或者插入作品集
func (r *WorksetsRepo) BulkUpsert(inputWorksets []*dbmodels.Workset) error {
	if len(inputWorksets) == 0 {
		return nil
	}

	// 使用 Clauses(clause.OnConflict{}) 实现批量 UPSERT
	// 冲突目标为 MoetranId，使冲突发生时，更新指定字段 Name 和 UpdatedAt
	if err := r.GetTable().
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{
				Name: "moetran_id", // 指定冲突字段为 moetran_id
			}},
			DoUpdates: clause.Assignments(map[string]any{ // 冲突时要更新的字段
				// EXCLUDED 表示冲突行的新值
				"name":       gorm.Expr("EXCLUDED.name"),
				"updated_at": gorm.Expr("EXCLUDED.updated_at"),
			}),
		}).
		Create(&inputWorksets).
		Error; err != nil {
		return fmt.Errorf("批量 Upsert Worksets 失败: %w", err)
	}

	return nil
}

// // SelectFullByTeamMoetranId 根据汉化组 ID 查询所有作品集
// // Full 代表 Workset 的完整信息
// func (r *WorksetsRepo) SelectFullByTeamId(teamId dbmodel.PrimaryKey) ([]*dbmodel.Workset, error) {
// 	var worksets []*dbmodel.Workset
// 	if err := r.GetTable().
// 		Where("team_id = ?", teamId).
// 		Find(&worksets).
// 		Error; err != nil {
// 		return nil, fmt.Errorf("查询 team_id %d 的作品集失败: %w", teamId, err)
// 	}

// 	return worksets, nil
// }
