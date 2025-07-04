package persistence

import (
	"fmt"
	"log/slog"
	"poplargrid/internal/shared/dbmodels"

	"gorm.io/gorm"
)

// 作品集 repo
type WorksetsRepo struct {
	// DbCtx 是数据库上下文
	DbCtx *gorm.DB
	// 日志器
	logger *slog.Logger // 如果需要日志记录，可以添加日志器
}

// NewWorksetsRepo 构造一个新的 WorksetsRepo 实例
func NewWorksetsRepo(db *gorm.DB, lgr *slog.Logger) *WorksetsRepo {
	return &WorksetsRepo{
		DbCtx:  db,
		logger: lgr,
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

	// 提取所有 MoetranId，用于批量查询
	moetranIds := make([]string, len(inputWorksets))
	for i, ws := range inputWorksets {
		moetranIds[i] = ws.MoetranId
	}

	// 批量查询已存在的 Workset 记录
	var existingWorksets []*dbmodels.Workset

	if err := r.GetTable().
		Select("id", "moetran_id", "updated_at").
		Where("moetran_id IN (?)", moetranIds).
		Find(&existingWorksets).
		Error; err != nil {
		return fmt.Errorf("批量查询 Worksets 失败: %w", err)
	}

	// 生成 MoetranId 为键的反查 map，方便查找
	idToSetMap := make(map[string]*dbmodels.Workset)
	for _, ws := range existingWorksets {
		idToSetMap[ws.MoetranId] = ws
	}

	// 区分需要创建和需要更新的记录
	toCreate := make([]*dbmodels.Workset, 0)
	toUpdate := make([]*dbmodels.Workset, 0)

	for _, inputWs := range inputWorksets {
		if existingWs, ok := idToSetMap[inputWs.MoetranId]; ok {
			// 记录已存在，检查是否需要更新
			if inputWs.UpdatedAt.After(existingWs.UpdatedAt) ||
				inputWs.Title != existingWs.Title {
				// 更新字段
				existingWs.Title = inputWs.Title
				existingWs.UpdatedAt = inputWs.UpdatedAt

				toUpdate = append(toUpdate, existingWs)
			}
			continue
		}
		// 记录不存在，需要创建
		toCreate = append(toCreate, inputWs)
	}

	// 批量创建新记录
	if len(toCreate) > 0 {
		r.logger.Info("批量创建 Workset 记录...", "count", len(toCreate))
		// 使用 CreateInBatches 进行批量插入
		if err := r.GetTable().CreateInBatches(toCreate, 100).Error; err != nil {
			return fmt.Errorf("批量创建 Worksets 失败: %w", err)
		}
	}

	// 批量更新现有记录 (GORM 的 Save 会根据主键更新，所以可以遍历更新)
	if len(toUpdate) > 0 {
		r.logger.Info("开始批量更新 Workset 记录...\n", "count", len(toUpdate))

		// 开启事务
		tx := r.DbCtx.Begin()
		if tx.Error != nil {
			return fmt.Errorf("开启事务失败: %w", tx.Error)
		}

		// 遍历更新每个 Workset
		for _, ws := range toUpdate {
			// Save 方法会根据主键更新所有字段
			if err := tx.Save(&ws).Error; err != nil {
				// 出现错误时回滚
				tx.Rollback()
				return fmt.Errorf("更新 Workset (MoetranId: %s) 失败: %w", ws.MoetranId, err)
			}
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("提交更新 Worksets 事务失败: %w", err)
		}
	}

	return nil
}
