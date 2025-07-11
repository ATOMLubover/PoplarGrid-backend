package dto

import (
	"fmt"
)

// ProjectListQueryParams 定义了获取项目列表的查询参数
type ProjectListQueryParams struct {
	PageSerial uint   `url:"pageSerial,default=1"` // 页码，默认值为 1
	PageSize   uint   `url:"pageSize,default=10"`  // 每页数量，默认值为 10
	Sort       string `url:"sort,default=id_desc"` // 排序方式，默认值为 id_desc，支持 id_desc | update_desc
	Status     string `url:"status,default=wip"`   // 项目状态，默认为 wip，支持 all | wip | fin
	TagId      string `url:"tag_id"`               // 特定标签，默认为空，表示不筛选
}

// Validate 方法进行额外校验或规范化
func (p *ProjectListQueryParams) Validate() error {
	// 校验 Sort 参数
	switch p.Sort {
	case "id_desc", "update_desc":
		// 有效值，无需处理
	default:
		// 其他情况视为错误
		return fmt.Errorf("无效的排序方式: %s", p.Sort)
	}

	// 校验 Status 参数
	switch p.Status {
	case "all", "translating", "translated", "prooving", "prooved", "lettering", "lettered", "published":
		// 有效值，无需处理
	default:
		// 其他情况视为错误
		return fmt.Errorf("无效的项目状态: %s", p.Status)
	}

	return nil
}
