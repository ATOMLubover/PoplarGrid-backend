package apiclient

// MoetranAppliCheckType 定义加入项目申请的处理方式
type MoetranAppliCheckType int

// 加入项目申请的处理方式
const (
	APPLI_NON_CHECK   MoetranAppliCheckType = iota // 允许所有人直接加入
	APPLI_ADMIN_CHECK                              // 仅允许管理员通过申请
)

// MoetranAllowApplyType 定义了项目的可申请范围
type MoetranAllowApplyType int

// 项目的可申请范围
const (
	ALLOW_NO_APPLI    MoetranAllowApplyType = iota // 禁止任何形式的申请
	ALLOW_ANY_APPLI                                // 允许任何形式的申请
	ALLOW_MEMBER_ONLY                              // 只允许本组成员申请
)

// MoetranRole 定义了龙译中的角色分工
type MoetranRole string

// 龙译角色的常量定义
const (
	ROLE_ADMIN       MoetranRole = "63d87c24b8bebd75ff934264" // 管理员
	ROLE_SUPERVISOR  MoetranRole = "63d87c24b8bebd75ff934265" // 监理
	ROLE_PROOFREADER MoetranRole = "63d87c24b8bebd75ff934266" // 校对
	ROLE_TRANSLATOR  MoetranRole = "63d87c24b8bebd75ff934267" // 翻译
	ROLE_EMBEDDER    MoetranRole = "63d87c24b8bebd75ff934268" // 嵌字
	ROLE_INTERN      MoetranRole = "63d87c24b8bebd75ff934269" // 实习翻译
)

// 语言代码
const (
	// East Asian Languages
	LangJapanese           string = "ja"    // 日语
	LangSimplifiedChinese  string = "zh-CN" // 简体中文
	LangTraditionalChinese string = "zh-TW" // 繁体中文
	LangKorean             string = "ko"    // 韩语

	// Western European Languages
	LangEnglish string = "en" // 英语
)

// normalErrorResponse 定义了尨译的标准错误响应格式
type normalErrorResponse struct {
	Code    int    `json:"code"`    // 错误代码
	Error   string `json:"error"`   // 错误信息
	Message string `json:"message"` // 错误详情
}

// CreateProjectParams 定义了创建项目的基本信息
type CreateProjectParams struct {
	MoetranAuth string // 龙译 JWT

	Title       string // 项目标题
	Description string // 项目简介

	MoetranProjSetId string // 龙译项目集 ID
	MoetranTeamId    string // 龙译团队 ID
	WorksetIndex     uint   // 工作集索引

	SourceLanguage  string   // 源语言
	TargetLanguages []string // 目标语言列表

	AllowApplyType       MoetranAllowApplyType // 允许申请类型
	ApplicationCheckType MoetranAppliCheckType // 申请审核类型

	DefaultRole MoetranRole // 默认角色 ID
}

// moetranCreateProjectRequest 定义了创建项目的请求结构体
type moetranCreateProjectRequest struct {
	Name                 string                `json:"name"`                   // 项目标题
	Intro                string                `json:"intro"`                  // 项目简介
	SourceLanguage       string                `json:"source_language"`        // 源语言
	TargetLanguages      []string              `json:"target_languages"`       // 目标语言列表
	AllowApplyType       MoetranAllowApplyType `json:"allow_apply_type"`       // 允许申请类型
	ApplicationCheckType MoetranAppliCheckType `json:"application_check_type"` // 申请审核类型
	DefaultRole          MoetranRole           `json:"default_role"`           // 默认角色 ID
	ProjectSet           string                `json:"project_set"`            // 项目集 ID
}

// moetranCreateProjectResponse 定义了创建项目的响应结构体
type moetranCreateProjectResponse struct {
	Message string `json:"message"` // 响应消息
	Project struct {
		Id string `json:"id"` // 项目 ID
	} `json:"project"` // 项目详情
}

// CreateProjectSetParams 定义了创建项目集的基本信息
type CreateProjectSetParams struct {
	MoetranAuth string // 龙译 JWT

	MoetranTeamId string // 龙译团队 ID
	Name          string // 项目集名称
}

// moetranCreateProjectSetRequest 定义了创建项目集的请求结构体
type moetranCreateProjectSetRequest struct {
	Name string `json:"name"` // 项目集名称
}

// moetranCreateProjectSetResponse 定义了创建项目集的响应结构体
type moetranCreateProjectSetResponse struct {
	Message    string `json:"message"` // 响应消息
	ProjectSet struct {
		Id string `json:"id"` // 项目集 ID
	} `json:"project_set"` // 项目集详情
}

// InviteMemberParams 定义了邀请成员加入特定项目的基本信息
type InviteMemberParams struct {
	MoetranAuth string // 龙译 JWT

	MoetranProjectId string      // 龙译项目 ID
	MoetranInviteeID string      // 龙译被邀请成员 ID
	InviteRole       MoetranRole // 邀请的角色 ID
}

// moetranInviteMemberRequest 定义了邀请成员的请求结构体
type moetranInviteMemberRequest struct {
	UserId  string `json:"user_id"` // 被邀请成员的 ID
	RoleId  string `json:"role"`    // 邀请的角色 ID
	Message string `json:"message"` // 邀请消息
}

// moetranInviteMemberResponse 定义了邀请成员的响应结构体
type moetranInviteMemberResponse struct {
	Message string `json:"message"` // 响应消息
}

// LoginParams 定义了登录龙译账号的参数
type LoginParams struct {
	Email       string // 龙译账号邮箱
	Password    string // 龙译账号密码
	Captcha     string // 验证码
	CaptchaInfo string // 验证码信息
}

// moetranLoginRequest 定义了登录请求的结构体
type moetranLoginRequest struct {
	Email       string `json:"email"`        // 账号邮箱
	Password    string `json:"password"`     // 账号密码
	Captcha     string `json:"captcha"`      // 验证码
	CaptchaInfo string `json:"captcha_info"` // 验证码信息
}

// moetranLoginResponse 定义了登录响应的结构体
type moetranLoginResponse struct {
	Token   string `json:"token"` // JWT
	Error   string `json:"error"`
	Message struct {
		Email    []string `json:"email"`    // 邮箱
		Password []string `json:"password"` // 密码
	} `json:"message"`
}

// RegisterParams 定义了注册龙译账号的参数
type RegisterParams struct {
	Email    string // 龙译账号邮箱
	Password string // 龙译账号密码
	Name     string // 龙译账号昵称
	VCode    string // 龙译发送的邮箱验证码
}

// moetranRegisterRequest 定义了注册请求的结构体
type moetranRegisterRequest struct {
	Email    string `json:"email"`    // 账号邮箱
	Password string `json:"password"` // 账号密码
	Name     string `json:"name"`     // 账号昵称
	VCode    string `json:"v_code"`   // 邮箱验证码
}

// moetranRegisterResponse 定义了注册响应的结构体
type moetranRegisterResponse struct {
	Message string `json:"message"` // 响应消息
	Token   string `json:"token"`   // JWT
}

// UserDTO 定义了获取用户信息的响应结构体
type UserDTO struct {
	ID   string `json:"id"`   // 用户 ID
	Name string `json:"name"` // 用户昵称
}

// UserInfo 定义了获取的用户信息的结构体
type UserInfo struct {
	Error    normalErrorResponse // 错误发生时的响应
	Response UserDTO             // 龙译用户信息响应
}

// TeamDTO 定义了获取汉化组信息的部分响应结构体
type TeamDTO struct {
	ID   string `json:"id"`   // 团队 ID
	Name string `json:"name"` // 团队名称
}

// UserTeamInfo 定义了获取用户团队信息的结构体
type UserTeamInfo struct {
	Error normalErrorResponse // 错误发生时的响应
	Teams []TeamDTO           // 用户所在的团队列表
}

// ProjectSetDTO 定义了获取项目集信息的部分响应结构体
type ProjectSetDTO struct {
	ID   string `json:"id"`   // 项目集 ID
	Name string `json:"name"` // 项目集名称
}

// ProjectSetInfo 定义了获取项目集信息的结构体
type ProjectSetInfo struct {
	Error normalErrorResponse // 错误发生时的响应
	Sets  []ProjectSetDTO     // 用户所在的项目集列表
}

// ProjectDTO 定义了获取项目信息的部分响应结构体
type ProjectDTO struct {
	ID    string `json:"id"`    // 项目 ID
	Name  string `json:"name"`  // 项目名称
	Intro string `json:"intro"` // 项目简介
}

// ProjectsInfo 定义了获取项目信息的结构体
type ProjectsInfo struct {
	Error    normalErrorResponse // 错误发生时的响应
	Projects []ProjectDTO        // 用户所在的项目列表
}
