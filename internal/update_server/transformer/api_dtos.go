package transformer

import (
	"fmt"
	"time"
)

// MoetranProjSet 是尨译关于作品集信息
type MoetranProjSet struct {
	CreateTime MoetranTime `json:"create_time"`
	EditTime   MoetranTime `json:"edit_time"`
	Default    bool        `json:"default"`
	Id         string      `json:"id"` // 龙译作品集 ID
	Intro      string      `json:"intro"`
	Name       string      `json:"name"`
}

// MoetranProj 是尨译关于项目/作品的信息
type MoetranProj struct {
	CreateTime MoetranTime `json:"create_time"` // 项目创建时间
	EditTime   MoetranTime `json:"edit_time"`   // 项目最后编辑时间
	Id         string      `json:"id"`          // 对应 Work.MoetranId

	Name  string `json:"name"`  // 包含 LegacyId，对应本地 Work.Title
	Intro string `json:"intro"` // 项目简介，暂时不确定能怎么对应到本地模型

	Status int `json:"status"` // 项目在第三方系统的状态，对应本地 Project.Status

	// 嵌套的所属作品集信息
	ProjectSet struct {
		Id string `json:"id"` // 对应Workset.MoetranId
	} `json:"project_set"`

	// 嵌套的所属汉化组信息
	Team struct {
		Id string `json:"id"` // 对应 Team.MoetranId
	} `json:"team"`
}

// ========= 以下是对尨译特殊的时间格式进行的反序列化辅助函数 =========

// 定义 MoetranTime 类型和其解析布局
// MOETRAN_TIME_LAYOUT 匹配 API 返回的 "YYYY-MM-DDTHH:mm:ss.SSSSSS" 格式
const MOETRAN_TIME_LAYOUT = "2006-01-02T15:04:05.000000"

// MoetranTime 用于正确解析 Moetran API 返回的时间字符串
type MoetranTime struct {
	time.Time
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口
// 当 JSON 解码器遇到 MoetranTime 字段时，会调用此方法
func (mt *MoetranTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	// 移除 JSON 字符串两端的双引号
	if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// 用定义的 CustomMoetranTimeLayout 来解析时间字符串
	// 如果解析失败，返回自定义错误信息
	parsedTime, err := time.Parse(MOETRAN_TIME_LAYOUT, s)
	if err != nil {
		return fmt.Errorf("无法解析时间字符串 '%s' 为格式 '%s': %w", s, MOETRAN_TIME_LAYOUT, err)
	}

	// 成功解析后，将解析结果赋值给 MoetranTime 的底层 time.Time
	mt.Time = parsedTime

	return nil
}

// MarshalJSON 实现了 json.Marshaler 接口
// 这将确保在将 MoetranTime 类型编码回 JSON 时也使用相同的格式。
func (mt MoetranTime) MarshalJSON() ([]byte, error) {
	// 将 time.Time 格式化为字符串，并用双引号包裹，以符合 JSON 字符串规范
	return []byte(fmt.Sprintf(`"%s"`, mt.Time.Format(MOETRAN_TIME_LAYOUT))), nil
}
