package handlers

// ErrorResponse 统一定义发生错误时的 JSON 响应格式
type ErrorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"` // 可选，在 service 层发生错误时提供详细信息
}

// SuccessResponse 统一定义成功响应的 JSON 格式
type SuccessResponse struct {
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"` // 可选，提供额外的成功信息
}

// // handlerPair 集成了 handler 和取消时的 respond 函数
// type handlerPair struct {
// 	Handler   func()
// 	Canceller func()
// }

// // thrustPreventer 用于去除同一客户端快速请求的抖动
// type thrustPreventer struct {
// 	// 锁保护下述字段的并发
// 	mtx sync.Mutex
// 	// 使用唯一标志符检查是否其正在处理
// 	record map[string]handlerPair
// 	keyMtx map[string]sync.Mutex
// }

// // DetachHandler 用于在指定时间没有被打断之后执行处理函数
// func (p *thrustPreventer) DetachHandler(
// 	identifier string, waitTime time.Duration,
// 	handler, cancelReponder func()) {
// 	// 先检查是否有同标识符的 handler 在等待运行
// 	p.mtx.Lock()

// 	pair, ok := p.record[identifier]
// 	if ok {
// 		// 如果当前有正在等待执行的 handler，则取消
// 		if pair.Canceller != nil {
// 			pair.Canceller()
// 		}
// 	}

// 	// 建立当前的执行函数必要信息
// 	cancalChan := make(chan struct{})
// 	canceller := func() { close(cancalChan) } // 通知当前 request 停止处理

// 	// 将当前标识符所在键值对占下
// 	p.record[identifier] = handlerPair{
// 		Handler:   handler,
// 		Canceller: canceller,
// 	}

// 	// 释放协程在后台待机准备执行
// 	go func() {
// 		select {
// 		case <-cancalChan:
// 			// 执行取消补偿函数
// 			cancelReponder()
// 		case <-time.After(waitTime):
// 			// 上锁，检查当前的处理函数是否是自己
// 			p.mtx.Lock()
// 			// 经过 waitTime，开始执行 handler
// 			handler()
// 			delete(p.record, identifier)
// 			p.mtx.Unlock()
// 		}
// 	}()

// 	// 在确保 cancelChan 有函数在监听后再释放锁
// 	p.mtx.Unlock()
// }
