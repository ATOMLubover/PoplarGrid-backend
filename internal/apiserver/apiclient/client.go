package apiclient

import "poplargrid/internal/apiserver/dtos"

// ApiClient 定义了与尨译 API 服务器交互的客户端接口
type ApiClient interface {
	// CreateProject 创建一个新的项目
	CreateProject(request *dtos.CreateProjectRequest, worksetIndex uint) error
}

// apiClientImpl 是 ApiClient 的实现
type apiClientImpl struct {
	// 这里可以添加需要的字段，例如 HTTP 客户端、API 基础 URL
}

// NewApiClient 创建一个新的 ApiClient 实例
func NewApiClient() ApiClient {
	return &apiClientImpl{
		// 初始化需要的字段
	}
}

// CreateProject 实现 ApiClient 接口的 CreateProject 方法
func (c *apiClientImpl) CreateProject(request *dtos.CreateProjectRequest, worksetIndex uint) error {
	// 这里实现与尨译 API 服务器的交互逻辑
	// 例如发送 HTTP POST 请求到特定的 API 端点
	// 并处理响应

	// 示例代码（需要根据实际 API 实现）：
	// resp, err := c.httpClient.Post(c.apiBaseURL+"/projects", "application/json", request)
	// if err != nil {
	//     return fmt.Errorf("failed to create project: %w", err)
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusCreated {
	//     return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }

	return nil // 返回 nil 表示成功
}
