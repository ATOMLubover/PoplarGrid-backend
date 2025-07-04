package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"poplargrid/internal/update_server/transformer"
	"time"
)

const (
	// 获取 project set 的 URL 模板
	// 需要 page, limit, word 查询参数
	PROJ_SET_API_FMT = "%s/v1/teams/%s/project-sets?%s"
	// 获取 project 的 URL 模板
	// 需要 page, limit, project_set, status, word 查询参数
	PROJ_API_FMT = "%s/v1/teams/%s/projects?%s"

	// 获取作品集时的分页大小
	PROJ_SET_PAGE_SIZE = 50
	// 获取作品时的分页大小
	PROJ_PAGE_SIZE = 15
)

var (
	// 默认的 User-Agent 头部信息，模拟浏览器请求
	DEFAULT_HEADERS = map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"Accept-Language": "zh-CN",
		// "Origin":             "https://moetran.com",
		// "Priority":           "u=1, i",
		// "Referer":            "https://moetran.com/",
		// "Sec-Ch-Ua":          `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`,
		// "Sec-Ch-Ua-Mobile":   "?0",
		// "Sec-Ch-Ua-Platform": `"Windows"`,
		// "Sec-Fetch-Dest":     "empty",
		// "Sec-Fetch-Mode":     "cors",
		// "Sec-Fetch-Site":     "same-site",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0",
	}
)

// ApiClient 封装了一个模拟的爬虫客户端
type ApiClient struct {
	// AuthToken token，由外部指定
	authTokenStr string
	// 基础 URL
	baseUrl string

	// 底层复用的 HTTP 客户端
	httpClient *http.Client
}

// NewApiClient 创建一个新的 ApiClient 实例
// authTokenStr 要单独调用函数在运行时修改
func NewApiClient(
	baseUrl string,
) *ApiClient {
	return &ApiClient{
		baseUrl: baseUrl,

		httpClient: &http.Client{
			Timeout: 10 * time.Second, // 设置超时时间为 10 秒
		},
	}
}

// 修改 authTokenStr
func (c *ApiClient) ModifyAuthToken(newAuthToken string) {
	c.authTokenStr = newAuthToken
}

// GetProjectSetUri 获取指定汉化组的作品集
func (c *ApiClient) GetProjectSetUri(teamMoetranId string) ([]transformer.MoetranProjSet, error) {
	allProjSets := make([]transformer.MoetranProjSet, 0)

	// 循环获取所有分页的作品集信息
	for page := 1; ; page++ {
		// 先构造查询 URL
		urlParams := url.Values{}
		urlParams.Set("page", fmt.Sprintf("%d", page))
		urlParams.Set("limit", fmt.Sprintf("%d", PROJ_SET_PAGE_SIZE))

		projSetUrl := fmt.Sprintf(PROJ_SET_API_FMT,
			c.baseUrl, teamMoetranId, urlParams.Encode())

		// 请求获取作品集信息
		res, err := c.sSendGetRequest(projSetUrl, c.authTokenStr)
		if err != nil {
			return nil, fmt.Errorf("请求获取作品集信息失败：%w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			// 此处暂时先直接读取请求体，且不处理重试
			bodyBytes, _ := io.ReadAll(res.Body)
			return nil, fmt.Errorf("获取作品集信息失败，HTTP 状态码：%d，响应体：%s",
				res.StatusCode, string(bodyBytes))
		}

		currProjSets := make([]transformer.MoetranProjSet, 0)

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&currProjSets); err != nil {
			return nil, fmt.Errorf("解析作品集信息失败：%w", err)
		}

		allProjSets = append(allProjSets, currProjSets...)

		if len(currProjSets) < PROJ_SET_PAGE_SIZE {
			// 读取到了最后一页，退出
			break
		}

		// 为了防止请求过快导致被限速，增加延时
		time.Sleep(200 * time.Millisecond)
	}

	return allProjSets, nil
}

// sSendGetRequest 用来辅助构造和发送 GET 请求
func (c *ApiClient) sSendGetRequest(url, authToken string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 GET 请求失败：%w", err)
	}

	// 设置默认的请求头
	for key, value := range DEFAULT_HEADERS {
		req.Header.Set(key, value)
	}

	// 添加 Authorization 头部
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	// 发送请求
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送 GET 请求失败：%w", err)
	}

	return res, nil
}
