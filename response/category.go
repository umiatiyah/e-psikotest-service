package response

type CategoryResponse struct {
	Name          string `json:"value"`
	MinScore      int    `json:"min_score"`
	Duration      int    `json:"duration"`
	LimitQuestion int    `json:"limit_question"`
	Message       BaseResponse
}
