package dto

type Search struct {
	Page       uint    `json:"page"`
	PageSize   uint    `json:"pageSize"`
	SortBy     string  `json:"sortBy"`
	SortOrder  string  `json:"sortOrder"`
	TextSearch string  `json:"textSearch"`
	Fields     []Field `json:"fields"`
}

type Field struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
