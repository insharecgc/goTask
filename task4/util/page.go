package util

// PageParam 分页请求参数
type PageParam struct {
	Page     int `form:"page" binding:"min=1"`     // 页码（默认1）
	PageSize int `form:"pageSize" binding:"min=1"` // 每页条数（默认10）
}

// PageResult 分页响应结果
type PageResult struct {
	List     interface{} `json:"list"`     // 当前页数据
	Total    int64       `json:"total"`    // 总条数
	Page     int         `json:"page"`     // 当前页码
	PageSize int         `json:"pageSize"` // 每页条数
	TotalPage int        `json:"totalPage"`// 总页数
}

// CalcPageResult 计算分页结果（填充totalPage等）
func CalcPageResult(list interface{}, total int64, page, pageSize int) *PageResult {
	totalPage := int(total + int64(pageSize) - 1) / pageSize // 向上取整
	return &PageResult{
		List:      list,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}
}
