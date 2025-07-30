package crawler

// ProjectSet 是尨译关于作品集的信息
type ProjectSet struct {
	Default bool   `json:"default"`
	ID      string `json:"id"` // 龙译作品集 ID
	Intro   string `json:"intro"`
	Name    string `json:"name"`
}

// Project 是尨译关于项目/作品的信息
type Project struct {
	ID   string `json:"id"`   // 对应 Project.MoetranID
	Name string `json:"name"` // 包含 LegacyID，对应本地 Project.Title
}

// User 是尨译关于用户信息的结构体
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// // ========= 以下是对尨译特殊的时间格式进行的反序列化辅助函数 =========

// // 定义 MoetranTime 类型和其解析布局
// // MOETRAN_TIME_LAYOUT 匹配 API 返回的 "YYYY-MM-DDTHH:mm:ss.SSSSSS" 格式
// const MOETRAN_TIME_LAYOUT = "2006-01-02T15:04:05.000000"

// // MoetranTime 用于正确解析 Moetran API 返回的时间字符串
// type MoetranTime struct {
// 	time.Time
// }

// // UnmarshalJSON 实现了 json.Unmarshaler 接口
// // 当 JSON 解码器遇到 MoetranTime 字段时，会调用此方法
// func (mt *MoetranTime) UnmarshalJSON(b []byte) (err error) {
// 	s := string(b)
// 	// 移除 JSON 字符串两端的双引号
// 	if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
// 		s = s[1 : len(s)-1]
// 	}

// 	// 用定义的 CustomMoetranTimeLayout 来解析时间字符串
// 	// 如果解析失败，返回自定义错误信息
// 	parsedTime, err := time.Parse(MOETRAN_TIME_LAYOUT, s)
// 	if err != nil {
// 		return fmt.Errorf("无法解析时间字符串 '%s' 为格式 '%s': %w", s, MOETRAN_TIME_LAYOUT, err)
// 	}

// 	// 成功解析后，将解析结果赋值给 MoetranTime 的底层 time.Time
// 	mt.Time = parsedTime

// 	return nil
// }

// // MarshalJSON 实现了 json.Marshaler 接口
// // 这将确保在将 MoetranTime 类型编码回 JSON 时也使用相同的格式。
// func (mt MoetranTime) MarshalJSON() ([]byte, error) {
// 	// 将 time.Time 格式化为字符串，并用双引号包裹，以符合 JSON 字符串规范
// 	return []byte(fmt.Sprintf(`"%s"`, mt.Time.Format(MOETRAN_TIME_LAYOUT))), nil
// }
