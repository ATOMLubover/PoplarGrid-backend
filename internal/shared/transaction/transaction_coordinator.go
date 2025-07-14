package transaction

import (
	"context"

	"gorm.io/gorm"
)

// TransactionCoordinator 封装事务处理回调函数
// 通常可以用于手动协调跨表事务
type TransactionCoordinator struct {
	dbCtx *gorm.DB
}

// NewTransactionCoordinator 创建一个新的事务协调器实例
func NewTransactionCoordinator(dbCtx *gorm.DB) *TransactionCoordinator {
	return &TransactionCoordinator{
		dbCtx: dbCtx,
	}
}

// RunInTransaction 执行一个事务，确保在回调函数中发生的错误会导致事务回滚
// 回调函数 handler 接受一个 *gorm.DB 参数，表示当前事务的数据库上下文
func (c *TransactionCoordinator) RunInTransaction(ctx context.Context, handler func(tx *gorm.DB) error) error {
	tx := c.dbCtx.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.WithContext(ctx).Transaction(handler); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
