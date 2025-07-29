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
	// GetProjectInfo 获取指定项目的详细信息，将响应体的 JSON 直接作为 string 返回
	GetProjectInfo(params *GetProjectParams) (string, error)

	// CreateProject 创建一个新的项目
	CreateProject(params *CreateProjectParams) (*moetranCreateProjectResponse, error)
	// CreateProjectSet 创建一个新的项目集
	CreateProjectSet(params *CreateProjectSetParams) (*moetranCreateProjectSetResponse, error)
	// InviteMemberToProject 邀请成员加入项目
	InviteMemberToProject(params *InviteMemberParams) (*moetranInviteMemberResponse, error)

	// Login 登录龙译账号，返回 JWT token
	Login(params *LoginParams) (string, error)
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

// GetProjectInfo 实现 ApiClient 接口的 GetProjectInfo 方法
func (c *apiClientImpl) GetProjectInfo(params *GetProjectParams) (string, error) {
	// 组装请求 URL
	url := fmt.Sprintf("%s/projects/%s",
		c.baseUrl, params.ProjectId)

	// 创建 HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error("创建龙译 get project 请求失败", slog.Any("error", err))
		return "", errors.New("无法构建龙译 get project 请求")
	}

	// 设置必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", params.MoetranAuth))

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 get project 请求失败", slog.Any("error", err))
		return "", errors.New("请求龙译 get project 失败")
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		c.logger.Error("龙译 get project 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
		)
		return "", fmt.Errorf("龙译 get project 请求失败，状态码: %d", res.StatusCode)
	}

	// 不解析响应体，直接读到一个 string 中返回
	var resBody bytes.Buffer
	if _, err := resBody.ReadFrom(res.Body); err != nil {
		c.logger.Error("读取龙译 get project 响应体失败", slog.Any("error", err))
		return "", errors.New("读取龙译 get project 响应体失败")
	}

	// 返回响应体的字符串形式
	return resBody.String(), nil
}

// CreateProject 实现 ApiClient 接口的 CreateProject 方法
func (c *apiClientImpl) CreateProject(params *CreateProjectParams) (*moetranCreateProjectResponse, error) {
	// 截断标题到 40 bytes，以满足龙译的限制
	title := truncateStringByRune(params.Title, 40)

	// 组装 POST 请求体
	body := moetranCreateProjectRequest{
		Name:  title,
		Intro: params.Description,

		SourceLanguage:  params.SourceLanguage,
		TargetLanguages: params.TargetLanguages,

		AllowApplyType:       params.AllowApplyType,
		ApplicationCheckType: params.ApplicationCheckType,

		DefaultRole: params.DefaultRole,
		ProjectSet:  params.MoetranProjSetId,
	}
	// 转换为 JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("转换龙译 create project 请求体失败", slog.Any("error", err))
		return nil, errors.New("转换龙译 create project 请求体失败")
	}

	// 组装请求 URL
	url := fmt.Sprintf("%s/teams/%s/projects",
		c.baseUrl, params.MoetranTeamId)

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
	if err != nil {
		c.logger.Error("创建龙译 create project 请求失败", slog.Any("error", err))
		return nil, errors.New("无法构建龙译 create project 请求")
	}

	// 设置一些必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", params.MoetranAuth))

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
	var response moetranCreateProjectResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.logger.Error("解析龙译 create project 响应失败", slog.Any("error", err))
		return nil, errors.New("解析龙译 create project 响应失败")
	}

	return &response, nil
}

// CreateProjectSet 实现 ApiClient 接口的 CreateProjectSet 方法
func (c *apiClientImpl) CreateProjectSet(params *CreateProjectSetParams) (*moetranCreateProjectSetResponse, error) {
	// 组装 POST 请求体
	body := moetranCreateProjectSetRequest{
		Name: params.Name,
	}

	// 转换为 JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("转换龙译 create project set 请求体失败", slog.Any("error", err))
		return nil, errors.New("转换龙译 create project set 请求体失败")
	}

	// 组装请求 URL
	url := fmt.Sprintf("%s/teams/%s/project-sets",
		c.baseUrl, params.MoetranTeamId)

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
	if err != nil {
		c.logger.Error("创建龙译 create project set 请求失败", slog.Any("error", err))
		return nil, errors.New("无法构建龙译 create project set 请求")
	}

	// 设置一些必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", params.MoetranAuth))

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 create project set 请求失败", slog.Any("error", err))
		return nil, errors.New("请求龙译 create project set 失败")
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		c.logger.Error("龙译 create project set 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
		)
		return nil, fmt.Errorf("龙译 create project set 请求失败，状态码: %d", res.StatusCode)
	}

	// 解析响应体
	var response moetranCreateProjectSetResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.logger.Error("解析龙译 create project set 响应失败", slog.Any("error", err))
		return nil, errors.New("解析龙译 create project set 响应失败")
	}

	return &response, nil
}

// InviteMemberToProject 实现 ApiClient 接口的 InviteMemberToProject 方法
func (c *apiClientImpl) InviteMemberToProject(params *InviteMemberParams) (*moetranInviteMemberResponse, error) {
	// 组装 POST 请求体
	body := moetranInviteMemberRequest{
		UserId:  params.MoetranInviteeID,
		RoleId:  string(params.InviteRole),
		Message: "", // 暂时直接留空
	}

	// 转换为 JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("转换龙译 invite member 请求体失败", slog.Any("error", err))
		return nil, errors.New("转换龙译 invite member 请求体失败")
	}

	// 组装请求 URL
	url := fmt.Sprintf("%s/projects/%s/invitations",
		c.baseUrl, params.MoetranProjectId)

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
	if err != nil {
		c.logger.Error("创建龙译 invite member 请求失败", slog.Any("error", err))
		return nil, errors.New("无法构建龙译 invite member 请求")
	}

	// 设置一些必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", params.MoetranAuth))

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 invite member 请求失败", slog.Any("error", err))
		return nil, errors.New("请求龙译 invite member 失败")
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		c.logger.Error("龙译 invite member 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
		)
		return nil, fmt.Errorf("龙译 invite member 请求失败，状态码: %d", res.StatusCode)
	}

	// 解析响应体
	var response moetranInviteMemberResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.logger.Error("解析龙译 invite member 响应失败", slog.Any("error", err))
		return nil, errors.New("解析龙译 invite member 响应失败")
	}

	return &response, nil
}

// Login 实现 ApiClient 接口的 Login 方法
func (c *apiClientImpl) Login(params *LoginParams) (string, error) {
	// 组装 POST 请求体
	body := moetranLoginRequest{
		Email:       params.Email,
		Password:    params.Password,
		Captcha:     params.Captcha,
		CaptchaInfo: params.CaptchaInfo,
	}

	// 转换为 JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("转换龙译 login 请求体失败", slog.Any("error", err))
		return "", errors.New("转换龙译 login 请求体失败")
	}

	// 组装请求 URL
	url := fmt.Sprintf("%s/login", c.baseUrl)

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
	if err != nil {
		c.logger.Error("创建龙译 login 请求失败", slog.Any("error", err))
		return "", errors.New("无法构建龙译 login 请求")
	}

	// 设置一些必要的请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 login 请求失败", slog.Any("error", err))
		return "", errors.New("请求龙译 login 失败")
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		c.logger.Error("龙译 login 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
		)
		return "", fmt.Errorf("龙译 login 请求失败，状态码: %d", res.StatusCode)
	}

	var response moetranLoginResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.logger.Error("解析龙译 login 响应失败", slog.Any("error", err))
		return "", errors.New("解析龙译 login 响应失败")
	}

	return response.Token, nil
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
