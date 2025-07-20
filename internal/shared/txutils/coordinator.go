package txutils

import (
	"context"
	"fmt"

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

// RunInTransaction 将 handler 封装在一个事务中执行，tx 为注入的事务上下文
// handler 函数可以返回一个错误和一个补偿函数，该补偿函数在出现错误并 *回滚之后* 被调用（不为 nil 时）
// 补偿函数无须处理 recovery，事务协调器会自动捕获 panic 并将其转换为错误
func (c *TransactionCoordinator) RunInTransaction(ctx context.Context, handler func(tx *gorm.DB) (error, func() error)) error {
	if handler == nil {
		// 先确保 handler 不为 nil
		return fmt.Errorf("handler 不能为 nil")
	}

	// 在事务回滚后执行的补偿函数，先记录并延迟执行时间
	var compensate func() error

	// 将 handler 函数适配为一个事务执行函数
	// 使得其签名满足 transaction handler，且捕获 panic
	handlerAdapter := func(tx *gorm.DB) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// 捕获 panic 并转换为错误
				err = fmt.Errorf("handler 中出现 panic: %v", r)
			}
		}()

		// 执行 handler 函数，并记录补偿函数和错误
		if err, compensate = handler(tx); err != nil {
			// 如果 handler 中出现错误，返回错误会自动导致事务回滚
			return
		}

		return
	}

	// 将 handler 函数适配为一个事务执行函数并运行
	if txErr := c.handle.WithContext(ctx).Transaction(handlerAdapter); txErr != nil {
		// 如果事务执行失败，延迟到此时执行补偿函数
		compErr := executeWithRecovery(compensate)
		// 将事务返回的和补偿函数返回的 error 一并返回
		return fmt.Errorf("事务执行失败: %w; 补偿函数执行中错误: %v", txErr, compErr)
	}

	return nil
}

// executeWithRecovery 安全地执行补偿函数，捕获可能的 panic
func executeWithRecovery(unsafeFunc func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// 捕获 panic 并转换为错误
			err = fmt.Errorf("补偿函数执行时出现 panic: %v", r)
		}
	}()

	if unsafeFunc != nil {
		err = unsafeFunc()
	}

	return
}
