package response

type CategoryResponse struct {
	Value         string `json:"value"`
	MinScore      int    `json:"min_score"`
	Duration      int    `json:"duration"`
	LimitQuestion int    `json:"limit_question"`
	Message       BaseResponse
}

type CategoryListResponse struct {
	Id            string `json:"id"`
	Value         string `json:"value"`
	MinScore      int    `json:"min_score"`
	Duration      int    `json:"duration"`
	LimitQuestion int    `json:"limit_question"`
}
