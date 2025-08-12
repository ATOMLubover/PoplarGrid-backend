package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	baseURL = "http://localhost:8081"
)

// 包级变量用于共享状态
var (
	jwtToken  string
	userInfo  *UserInfo
	memberID  int
	teamID    int
	worksetID int
	projectID int
	appID     int
	inviteID  int
)

// 初始化测试
func TestMain(m *testing.M) {
	// 1. 初始化测试环境
	setup()

	// 2. 运行所有测试
	exitCode := m.Run()

	// 3. 清理测试数据
	teardown()

	// 4. 退出测试
	os.Exit(exitCode)
}

func setup() {
	fmt.Println("=== 测试初始化 ===")

	// 读取环境变量
	email := os.Getenv("POPLAR_EMAIL")
	password := os.Getenv("POPLAR_PASSWORD")
	captcha := os.Getenv("POPLAR_CAPTCHA")
	captchaInfo := os.Getenv("POPLAR_CAPTCHA_INFO")

	if email == "" || password == "" {
		log.Fatal("请设置环境变量 POPLAR_EMAIL 和 POPLAR_PASSWORD")
	}
	if captcha == "" {
		captcha = "000000" // 测试环境下默认验证码
	}
	if captchaInfo == "" {
		captchaInfo = "test-info"
	}

	// 执行绑定和登录
	testBindAccount(email, password, captcha, captchaInfo)
	testLogin(email, password, captcha, captchaInfo)

	fmt.Println("=== 初始化完成 ===")
}

func teardown() {
	fmt.Println("=== 测试清理 ===")

	if projectID != 0 {
		err := deleteProject(projectID)
		if err != nil {
			fmt.Printf("清理项目失败: %v\n", err)
		} else {
			fmt.Printf("已删除项目 %d\n", projectID)
		}
	}

	if worksetID != 0 {
		// 实际API中没有删除作品集的接口，这里只做演示
		fmt.Printf("清理作品集 %d\n", worksetID)
	}

	fmt.Println("=== 清理完成 ===")
}

// ==== 认证测试 ====

func testBindAccount(email, password, captcha, captchaInfo string) {
	fmt.Println("- 测试绑定龙译账号")
	params := BindParams{
		Email:       email,
		Password:    password,
		Captcha:     captcha,
		CaptchaInfo: captchaInfo,
	}

	err := bindAccount(params)
	if err != nil {
		log.Fatalf("绑定失败: %v", err)
	}
}

func testLogin(email, password, captcha, captchaInfo string) {
	fmt.Println("- 测试登录")
	params := LoginParams{
		Email:       email,
		Password:    password,
		Captcha:     captcha,
		CaptchaInfo: captchaInfo,
	}

	err := login(params)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
}

// ==== 用户信息测试 ====

func TestGetUserInfo(t *testing.T) {
	fmt.Println("- 测试获取用户信息")
	user, err := getUserInfo()
	if err != nil {
		t.Fatalf("获取用户信息失败: %v", err)
	}

	if len(user.Members) == 0 {
		t.Fatal("用户没有加入任何汉化组")
	}

	memberID = user.Members[0].ID
	teamID = user.Members[0].Team.ID
	userInfo = user

	fmt.Printf("用户ID: %d, 昵称: %s\n", user.ID, user.Nickname)
	fmt.Printf("团队成员ID: %d, 团队ID: %d\n", memberID, teamID)
}

// ==== 作品集测试 ====

func TestCreateWorkset(t *testing.T) {
	fmt.Println("- 测试创建作品集")
	req := CreateWorksetRequest{
		Name:        "测试作品集-" + randomSuffix(),
		Description: "API测试创建的作品集",
		TeamID:      teamID,
	}

	id, err := createWorkset(req)
	if err != nil {
		t.Fatalf("创建作品集失败: %v", err)
	}

	worksetID = id
	fmt.Printf("创建作品集成功，ID: %d\n", worksetID)
}

// ==== 项目测试 ====

func TestCreateProject(t *testing.T) {
	fmt.Println("- 测试创建项目")
	req := CreateProjectRequest{
		Title:             "API测试项目-" + randomSuffix(),
		WorksetID:         worksetID,
		ApplicantMemberID: memberID,
		Description:       "API测试创建的项目",
		AllowAutoJoin:     true,
	}

	id, err := createProject(req)
	if err != nil {
		t.Fatalf("创建项目失败: %v", err)
	}

	projectID = id
	fmt.Printf("创建项目成功，ID: %d\n", projectID)
}

func TestUpdateProject(t *testing.T) {
	fmt.Println("- 测试更新项目")
	req := UpdateProjectRequest{
		Title:            "更新后的测试项目-" + randomSuffix(),
		Description:      "已更新的项目描述",
		Status:           2,
		OperatorMemberID: memberID,
	}

	err := updateProject(projectID, req)
	if err != nil {
		t.Fatalf("更新项目失败: %v", err)
	}

	fmt.Println("项目更新成功")
}

func TestGetProjectDetail(t *testing.T) {
	fmt.Println("- 测试获取项目详情")
	project, err := getProjectDetail(projectID)
	if err != nil {
		t.Fatalf("获取项目详情失败: %v", err)
	}

	fmt.Printf("项目标题: %s\n描述: %s\n状态: %d\n",
		project.Title, project.Description, project.Status)
}

// ==== 协作流程测试 ====

func TestCreateApplication(t *testing.T) {
	fmt.Println("- 测试创建申请")
	req := CreateApplicationRequest{
		ApplicantMemberID: memberID,
		TargetLaborMask:   4, // 美工
		TargetProjectID:   projectID,
	}

	id, err := createApplication(req)
	if err != nil {
		t.Fatalf("创建申请失败: %v", err)
	}

	appID = id
	fmt.Printf("创建申请成功，ID: %d\n", appID)
}

func TestProcessApplication(t *testing.T) {
	fmt.Println("- 测试处理申请")
	req := ProcessApplicationRequest{
		Accept:            true,
		ApplicationID:     appID,
		ProcessorMemberID: memberID,
	}

	err := processApplication(req)
	if err != nil {
		t.Fatalf("处理申请失败: %v", err)
	}

	fmt.Println("处理申请成功")
}

func TestCreateInvitation(t *testing.T) {
	fmt.Println("- 测试创建邀请")
	req := CreateInvitationRequest{
		InvitorMemberID: memberID,
		InviteeMemberID: memberID,
		TargetLaborMask: 8, // 翻译
		TargetProjectID: projectID,
	}

	id, err := createInvitation(req)
	if err != nil {
		t.Fatalf("创建邀请失败: %v", err)
	}

	inviteID = id
	fmt.Printf("创建邀请成功，ID: %d\n", inviteID)
}

func TestProcessInvitation(t *testing.T) {
	fmt.Println("- 测试处理邀请")
	req := ProcessInvitationRequest{
		Accept:            true,
		InvitationID:      inviteID,
		ProcessorMemberID: memberID,
	}

	err := processInvitation(req)
	if err != nil {
		t.Fatalf("处理邀请失败: %v", err)
	}

	fmt.Println("处理邀请成功")
}

// ==== 统计信息测试 ====

func TestGetWorksetStats(t *testing.T) {
	fmt.Println("- 测试获取作品集统计信息")
	stats, err := getWorksetStats(worksetID)
	if err != nil {
		t.Fatalf("获取统计信息失败: %v", err)
	}

	fmt.Printf("作品集 %d 统计信息:\n", worksetID)
	fmt.Printf("总项目数: %d\n翻译中: %d\n已校对: %d\n已嵌字: %d\n已发布: %d\n",
		stats.TotalProjectCount, stats.OnTranslatingCount, stats.ProofreadCount,
		stats.LetteredCount, stats.PublishedCount)
}

// ==== 系统功能测试 ====

func TestTriggerAutoUpdate(t *testing.T) {
	fmt.Println("- 测试触发爬虫自动更新")
	err := triggerAutoUpdate()
	if err != nil {
		t.Fatalf("触发自动更新失败: %v", err)
	}

	fmt.Println("已触发爬虫自动更新")
}

// ==== 辅助函数 ====

func randomSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%10000)
}

// ==== HTTP请求辅助函数 ====

func sendRequest(method, url string, body interface{}, response interface{}) error {
	var reqBody []byte
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = jsonData
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	if jwtToken != "" {
		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return fmt.Errorf("解析响应失败: %v", err)
		}
	}
	return nil
}

// ==== API函数 ====

func bindAccount(params BindParams) error {
	url := baseURL + "/auth/bind"
	return sendRequest("POST", url, params, nil)
}

func login(params LoginParams) error {
	url := baseURL + "/auth/login"
	type Response struct {
		Token string `json:"token"`
	}
	var resp Response
	if err := sendRequest("POST", url, params, &resp); err != nil {
		return err
	}
	jwtToken = resp.Token
	return nil
}

func getUserInfo() (*UserInfo, error) {
	url := baseURL + "/api/users/me"
	var user UserInfo
	if err := sendRequest("GET", url, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func createWorkset(req CreateWorksetRequest) (int, error) {
	url := baseURL + "/api/worksets"
	type Response struct {
		ID int `json:"id"`
	}
	var resp Response
	if err := sendRequest("POST", url, req, &resp); err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func createProject(req CreateProjectRequest) (int, error) {
	url := baseURL + "/api/projects"
	type Response struct {
		ID int `json:"id"`
	}
	var resp Response
	if err := sendRequest("POST", url, req, &resp); err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func updateProject(id int, req UpdateProjectRequest) error {
	url := fmt.Sprintf("%s/api/projects/%d", baseURL, id)
	return sendRequest("PATCH", url, req, nil)
}

func getProjectDetail(id int) (*ProjectInfo, error) {
	url := fmt.Sprintf("%s/api/projects/%d", baseURL, id)
	var project ProjectInfo
	if err := sendRequest("GET", url, nil, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func createApplication(req CreateApplicationRequest) (int, error) {
	url := baseURL + "/api/applications"
	type Response struct {
		ID int `json:"id"`
	}
	var resp Response
	if err := sendRequest("POST", url, req, &resp); err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func processApplication(req ProcessApplicationRequest) error {
	url := fmt.Sprintf("%s/api/applications/%d", baseURL, req.ApplicationID)
	return sendRequest("PUT", url, req, nil)
}

func createInvitation(req CreateInvitationRequest) (int, error) {
	url := baseURL + "/api/invitations"
	type Response struct {
		ID int `json:"id"`
	}
	var resp Response
	if err := sendRequest("POST", url, req, &resp); err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func processInvitation(req ProcessInvitationRequest) error {
	url := fmt.Sprintf("%s/api/invitations/%d", baseURL, req.InvitationID)
	return sendRequest("PUT", url, req, nil)
}

func getWorksetStats(id int) (*WorksetStats, error) {
	url := fmt.Sprintf("%s/api/worksets/%d/stats", baseURL, id)
	var stats WorksetStats
	if err := sendRequest("GET", url, nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func triggerAutoUpdate() error {
	url := baseURL + "/api/crawler/auto-update-all"
	return sendRequest("POST", url, nil, nil)
}

func deleteProject(id int) error {
	url := fmt.Sprintf("%s/api/projects/%d", baseURL, id)
	return sendRequest("DELETE", url, nil, nil)
}

// ==== 数据结构 ====

type BindParams struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Captcha     string `json:"captcha"`
	CaptchaInfo string `json:"captcha_info"`
}

type LoginParams struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Captcha     string `json:"captcha"`
	CaptchaInfo string `json:"captcha_info"`
}

type CreateWorksetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TeamID      int    `json:"team_id"`
}

type CreateProjectRequest struct {
	Title             string `json:"title"`
	WorksetID         int    `json:"workset_id"`
	ApplicantMemberID int    `json:"applicant_member_id,omitempty"`
	Description       string `json:"description,omitempty"`
	AllowAutoJoin     bool   `json:"allow_auto_join,omitempty"`
	IsHidden          bool   `json:"is_hidden,omitempty"`
}

type UpdateProjectRequest struct {
	Title            string `json:"title"`
	Description      string `json:"description,omitempty"`
	Status           int    `json:"status,omitempty"`
	OperatorMemberID int    `json:"operator_member_id"`
}

type CreateApplicationRequest struct {
	ApplicantMemberID int `json:"applicant_member_id"`
	TargetLaborMask   int `json:"target_labor_mask"`
	TargetProjectID   int `json:"target_project_id"`
}

type CreateInvitationRequest struct {
	InvitorMemberID int `json:"invitor_member_id"`
	InviteeMemberID int `json:"invitee_member_id,omitempty"`
	TargetLaborMask int `json:"target_labor_mask"`
	TargetProjectID int `json:"target_project_id"`
}

type ProcessApplicationRequest struct {
	Accept            bool `json:"accept"`
	ApplicationID     int  `json:"application_id"`
	ProcessorMemberID int  `json:"processor_member_id"`
}

type ProcessInvitationRequest struct {
	Accept            bool `json:"accept"`
	InvitationID      int  `json:"invitation_id"`
	ProcessorMemberID int  `json:"processor_member_id"`
}

type UserInfo struct {
	ID         int          `json:"id"`
	Nickname   string       `json:"nickname"`
	Email      string       `json:"email"`
	IsAdmin    bool         `json:"is_admin"`
	MoetranID  string       `json:"moetran_id"`
	MoetranJWT string       `json:"moetran_jwt"`
	QQNumber   int          `json:"qq_number"`
	Remark     string       `json:"remark"`
	Members    []MemberInfo `json:"members"`
}

type MemberInfo struct {
	ID   int      `json:"id"`
	Role int      `json:"role"`
	Team TeamInfo `json:"team"`
	User UserInfo `json:"user"`
}

type TeamInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MoetranID   string `json:"moetran_id"`
}

type ProjectInfo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	WorksetID   int    `json:"workset_id"`
}

type WorksetStats struct {
	TotalProjectCount    int `json:"total_project_count"`
	NotTranslatingCount  int `json:"not_translating_count"`
	NotProofreadingCount int `json:"not_proofreading_count"`
	NotLetteringCount    int `json:"not_lettering_count"`
	NotReviewingCount    int `json:"not_reviewing_count"`
	OnTranslatingCount   int `json:"on_translating_count"`
	OnProofreadingCount  int `json:"on_proofreading_count"`
	OnLetteringCount     int `json:"on_lettering_count"`
	OnReviewingCount     int `json:"on_reviewing_count"`
	TranslatedCount      int `json:"translated_count"`
	ProofreadCount       int `json:"proofread_count"`
	LetteredCount        int `json:"lettered_count"`
	ReviewedCount        int `json:"reviewed_count"`
	PublishedCount       int `json:"published_count"`
}
