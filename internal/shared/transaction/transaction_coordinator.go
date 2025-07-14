package transaction

import (
	"context"

	"gorm.io/gorm"
)

// TransactionCoordinator 封装事务处理回调函数
// 通常可以用于手动协调跨表事务
type TransactionCoordinator struct {
	handle *gorm.DB
}

// NewTransactionCoordinator 创建一个新的事务协调器实例
func NewTransactionCoordinator(handle *gorm.DB) *TransactionCoordinator {
	return &TransactionCoordinator{
		handle: handle,
	}
}

// RunInTransaction 执行一个事务，确保在回调函数中发生的错误会导致事务回滚
// 回调函数 handler 接受一个 *gorm.DB 参数，表示当前事务的数据库上下文
// handler 函数可以返回一个错误和一个补偿函数，清理函数在出现错误时被调用
func (c *TransactionCoordinator) RunInTransaction(ctx context.Context, handler func(tx *gorm.DB) (error, func())) error {
	tx := c.handle.Begin()

	// 防止出现 panic 时干崩服务器，且事务未被回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 包装 handler 函数，使其能够在事务中执行
	txHandler := func(tx *gorm.DB) error {
		err, cleanup := handler(tx)
		if err != nil && cleanup != nil {
			// 当出现错误，且 cleanup 函数不为 nil 时，执行清理操作
			defer cleanup()
		}

		// 如果 handler 返回错误，上传到事务中
		return err
	}

	if err := tx.WithContext(ctx).Transaction(txHandler); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
