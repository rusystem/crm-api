package domain

type Param struct {
	Limit     int64  `json:"limit"`
	Offset    int64  `json:"offset"`
	CompanyId int64  `json:"company_id"`
	Query     string `json:"query"`
	Sort      string `json:"sort"`
	SortField string `json:"sort_field"`
}
