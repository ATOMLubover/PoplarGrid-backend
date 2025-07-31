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
	CreateProject(params *CreateProjectParams) (*moetranCreateProjectResponse, error)
	// CreateProjectSet 创建一个新的项目集
	CreateProjectSet(params *CreateProjectSetParams) (*moetranCreateProjectSetResponse, error)
	// InviteMemberToProject 邀请成员加入项目
	InviteMemberToProject(params *InviteMemberParams) (*moetranInviteMemberResponse, error)

	// Login 登录龙译账号，返回 JWT 或者错误信息
	Login(params *LoginParams) (string, error)
	// Register 注册龙译账号，返回 JWT 或者错误信息
	Register(params *RegisterParams) (string, error)

	// GetProjects 获取指定项目集下的部分项目
	GetProjects(teamID, projectSetId string, page, limit int, moetranAuth string) (*ProjectsInfo, error)
	// GetUserInfo 获取指定用户的详细信息
	GetUserInfo(moetranAuth string) (*UserInfo, error)
	// GetUserTeams 获取指定用户所在的团队列表
	GetUserTeams(moetranAuth string) (*UserTeamInfo, error)
	// GetTeamProjectSets 获取指定汉化组的项目集列表
	GetTeamProjectSets(teamId string, moetranAuth string) (*ProjectSetInfo, error)
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
	url := fmt.Sprintf("%s/user/token", c.baseUrl)

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

	var response moetranLoginResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.logger.Error("解析龙译 login 响应失败", slog.Any("error", err))
		return "", errors.New("解析龙译 login 响应失败")
	}

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		c.logger.Error("龙译 login 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url))
		return "", fmt.Errorf("龙译 login 请求失败，状态码: %d，龙译错误：%v",
			res.StatusCode, response.Message)
	}

	return response.Token, nil
}

// Register 实现 ApiClient 接口的 Register 方法
func (c *apiClientImpl) Register(params *RegisterParams) (string, error) {
	// 组装 POST 请求体
	body := moetranRegisterRequest{
		Email:    params.Email,
		Password: params.Password,
		Name:     params.Name,
		VCode:    params.VCode,
	}

	// 转换为 JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("转换龙译 register 请求体失败", slog.Any("error", err))
		return "", errors.New("转换龙译 register 请求体失败")
	}

	// 组装请求 URL
	url := fmt.Sprintf("%s/users", c.baseUrl)

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJson))
	if err != nil {
		c.logger.Error("创建龙译 register 请求失败", slog.Any("error", err))
		return "", errors.New("无法构建龙译 register 请求")
	}

	// 设置一些必要的请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 register 请求失败", slog.Any("error", err))
		return "", errors.New("请求龙译 register 失败")
	}
	defer res.Body.Close()

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		c.logger.Error("龙译 register 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
		)
		return "", fmt.Errorf("龙译 register 请求失败，状态码: %d", res.StatusCode)
	}

	var response moetranRegisterResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.logger.Error("解析龙译 register 响应失败", slog.Any("error", err))
		return "", errors.New("解析龙译 register 响应失败")
	}

	return response.Token, nil
}

// GetUserInfo 实现 ApiClient 接口的 GetUserInfo 方法
func (c *apiClientImpl) GetUserInfo(moetranAuth string) (*UserInfo, error) {
	// 组装请求 URL
	url := fmt.Sprintf("%s/user/info", c.baseUrl)

	// 创建 HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error("创建龙译 get user info 请求失败", slog.Any("error", err))
		return nil, errors.New("无法构建龙译 get user info 请求")
	}

	// 设置必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", moetranAuth))

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 get user info 请求失败", slog.Any("error", err))
		return nil, errors.New("请求龙译获取用户信息失败")
	}
	defer res.Body.Close()

	userInfo := &UserInfo{}

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		// 此时尝试读取为 error 信息
		if err := json.NewDecoder(res.Body).Decode(&userInfo.Error); err != nil {
			c.logger.Error("解析龙译 get user info 错误响应失败", slog.Any("error", err))
			return nil, errors.New("解析龙译错误响应失败")
		}

		c.logger.Error("龙译 get user info 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
			slog.String("error_message", userInfo.Error.Message))
		return userInfo, fmt.Errorf("龙译获取用户信息请求失败，状态码: %d", res.StatusCode)
	}

	// 解析正常的响应体
	if err := json.NewDecoder(res.Body).Decode(&userInfo.Response); err != nil {
		c.logger.Error("解析龙译 get user info 响应失败", slog.Any("error", err))
		return nil, errors.New("解析龙译获取用户信息响应失败")
	}

	return userInfo, nil
}

// GetUserTeams 实现 ApiClient 接口的 GetUserTeams 方法
func (c *apiClientImpl) GetUserTeams(moetranAuth string) (*UserTeamInfo, error) {
	// 组装请求 URL
	url := fmt.Sprintf("%s/user/teams", c.baseUrl)

	// 创建 HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error("创建龙译 get user teams 请求失败", slog.Any("error", err))
		return nil, errors.New("无法构建龙译 get user teams 请求")
	}

	// 设置必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", moetranAuth))

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 get user teams 请求失败", slog.Any("error", err))
		return nil, errors.New("请求龙译获取用户团队信息失败")
	}
	defer res.Body.Close()

	userTeamInfo := &UserTeamInfo{}

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		// 此时尝试读取为 error 信息
		if err := json.NewDecoder(res.Body).Decode(&userTeamInfo.Error); err != nil {
			c.logger.Error("解析龙译 get user teams 错误响应失败", slog.Any("error", err))
			return nil, errors.New("解析龙译错误响应失败")
		}

		c.logger.Error("龙译 get user teams 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
			slog.String("error_message", userTeamInfo.Error.Message))
		return userTeamInfo, fmt.Errorf("龙译获取用户团队信息请求失败，状态码: %d", res.StatusCode)
	}

	// 解析正常的响应体
	if err := json.NewDecoder(res.Body).Decode(&userTeamInfo.Teams); err != nil {
		c.logger.Error("解析龙译 get user teams 响应失败", slog.Any("error", err))
		return nil, errors.New("解析龙译获取用户团队信息响应失败")
	}

	return userTeamInfo, nil
}

// GetTeamProjectSets 实现 ApiClient 接口的 GetTeamProjectSets 方法
func (c *apiClientImpl) GetTeamProjectSets(teamId string, moetranAuth string) (*ProjectSetInfo, error) {
	// 组装请求 URL
	url := fmt.Sprintf("%s/teams/%s/project-sets", c.baseUrl, teamId)

	// 创建 HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error("创建龙译 get team project sets 请求失败", slog.Any("error", err))
		return nil, errors.New("无法构建龙译 get team project sets 请求")
	}

	// 设置必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", moetranAuth))

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 get team project sets 请求失败", slog.Any("error", err))
		return nil, errors.New("请求龙译获取团队项目集信息失败")
	}
	defer res.Body.Close()

	projectSetInfo := &ProjectSetInfo{}

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		// 此时尝试读取为 error 信息
		if err := json.NewDecoder(res.Body).Decode(&projectSetInfo.Error); err != nil {
			c.logger.Error("解析龙译 get team project sets 错误响应失败", slog.Any("error", err))
			return nil, errors.New("解析龙译错误响应失败")
		}

		c.logger.Error("龙译 get team project sets 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
			slog.String("error_message", projectSetInfo.Error.Message))
		return projectSetInfo, fmt.Errorf("龙译获取团队项目集信息请求失败，状态码: %d", res.StatusCode)
	}

	// 解析正常的响应体
	if err := json.NewDecoder(res.Body).Decode(&projectSetInfo.Sets); err != nil {
		c.logger.Error("解析龙译 get team project sets 响应失败", slog.Any("error", err))
		return nil, errors.New("解析龙译获取团队项目集信息响应失败")
	}

	return projectSetInfo, nil
}

// GetProjects 实现 ApiClient 接口的 GetProjects 方法
func (c *apiClientImpl) GetProjects(teamID, projectSetID string, page, limit int, moetranAuth string) (*ProjectsInfo, error) {
	// 组装请求 URL
	url := fmt.Sprintf("%s/teams/%s/projects?page=%d&limit=%d&project_set=%s",
		c.baseUrl, teamID, page, limit, projectSetID)

	// 创建 HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Error("创建龙译 get projects 请求失败", slog.Any("error", err))
		return nil, errors.New("无法构建龙译 get projects 请求")
	}

	// 设置必要的请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", moetranAuth))

	// 发送请求
	res, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("发送龙译 get projects 请求失败", slog.Any("error", err))
		return nil, errors.New("请求龙译获取项目列表失败")
	}
	defer res.Body.Close()

	projectsInfo := &ProjectsInfo{}

	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		// 此时尝试读取为 error 信息
		if err := json.NewDecoder(res.Body).Decode(&projectsInfo.Error); err != nil {
			c.logger.Error("解析龙译 get projects 错误响应失败", slog.Any("error", err))
			return nil, errors.New("解析龙译错误响应失败")
		}

		c.logger.Error("龙译 get projects 请求失败",
			slog.Int("status_code", res.StatusCode),
			slog.String("url", url),
			slog.Any("error_message", projectsInfo.Error))
		return projectsInfo, fmt.Errorf("龙译获取项目列表请求失败，状态码: %d", res.StatusCode)
	}

	// 解析正常的响应体
	if err := json.NewDecoder(res.Body).Decode(&projectsInfo.Projects); err != nil {
		c.logger.Error("解析龙译 get projects 响应失败", slog.Any("error", err))
		return nil, errors.New("解析龙译获取项目列表响应失败")
	}

	return projectsInfo, nil
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
