package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"unicode/utf8"
)

// ApiClient 定义了与尨译 API 服务器交互的客户端接口
type ApiClient interface {
	// CreateProject 创建一个新的项目
	CreateProject(info *CreateProjectInfo) (*MoetranCreateProjectResponse, error)
}

// NewApiClient 创建一个新的 ApiClient 实例
// 目前的 baseUrl 应该是 https://api.moetran.com/v1
func NewApiClient(
	baseUrl string,
	logger slog.Logger,
) ApiClient {
	// 这里可以根据需要返回具体的实现
	return &apiClientImpl{
		client:  createHttpClient(),
		baseUrl: baseUrl,
		logger:  logger,
	}
}

// apiClientImpl 是 ApiClient 的实现
type apiClientImpl struct {
	// 这里可以添加需要的字段，例如 HTTP 客户端、API 基础 URL

	client  *http.Client
	baseUrl string

	logger slog.Logger
}

// CreateProject 实现 ApiClient 接口的 CreateProject 方法
func (c *apiClientImpl) CreateProject(info *CreateProjectInfo) (*MoetranCreateProjectResponse, error) {
	// 截断标题到 40 bytes，以满足龙译的限制
	title := truncateStringByRune(info.Title, 40)

	// 组装 POST 请求体
	body := MoetranCreateProjectRequest{
		Name:  title,
		Intro: info.Description,

		SourceLanguage:  info.SourceLanguage,
		TargetLanguages: info.TargetLanguages,

		AllowApplyType:       info.AllowApplyType,
		ApplicationCheckType: info.ApplicationCheckType,

		DefaultRole: info.DefaultRole,
		ProjectSet:  info.MoetranProjSetId,
	}
	// 转换为 JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("转换龙译 create project 请求体失败", slog.Any("error", err))
		return nil, errors.New("转换龙译 create project 请求体失败")
	}

	// 组装请求 URL
	url := fmt.Sprintf("%s/teams/%s/projects",
		c.baseUrl, info.MoetranTeamId)

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
	if err != nil {
		c.logger.Error("创建龙译 create project 请求失败", slog.Any("error", err))
		return nil, errors.New("无法构建龙译 create project 请求")
	}

	// 设置一些必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", info.MoetranAuth))

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 create project 请求失败", slog.Any("error", err))
		return nil, errors.New("请求龙译 create project 失败")
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		c.logger.Error("龙译 create project 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
		)
		return nil, fmt.Errorf("龙译 create project 请求失败，状态码: %d", res.StatusCode)
	}

	// 解析响应体
	var response MoetranCreateProjectResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.logger.Error("解析龙译 create project 响应失败", slog.Any("error", err))
		return nil, errors.New("解析龙译 create project 响应失败")
	}

	return &response, nil
}

// =========== 辅助函数 ===========

// createHttpClient 创建一个新的 HTTP 客户端
func createHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

// truncateStringByRune 截断字符串到指定的长度（以 字符/char 为单位）
func truncateStringByRune(s string, lenInBytes int) string {
	if len(s) <= lenInBytes {
		return s
	}

	// 先尝试截取到指定长度的字节数
	// 此时可能会截断在一个 UTF-8 字符的中间
	truncated := []byte(s[:lenInBytes])

	// 然后试探检查是否是合法的 UTF-8 字符串
	for {
		// utf8.DecodeLastRune 尝试解码 truncated 中最后一个完整的 UTF-8 字符
		// 其它返回解码出的 rune 和该 rune 所占的字节数
		r, size := utf8.DecodeLastRune(truncated)

		if size != 0 && r != utf8.RuneError {
			// 如果解码成功，并且不是 RuneError，可以跳出试探循环
			break
		}

		// 至此，说明最后一个字符不是合法的 UTF-8 字符
		if len(truncated) > 0 {
			// 如果 truncated 还不为空，继续剪除最后一个 char
			truncated = truncated[:len(truncated)-1]
			continue
		}

		// 如果 truncated 试探截取到空，直接 return 空字符串
		return ""
	}

	return string(truncated)
}
