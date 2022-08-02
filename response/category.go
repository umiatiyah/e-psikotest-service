package response

type CategoryResponse struct {
	Name          string `json:"value"`
	MinScore      int    `json:"min_score"`
	Duration      int    `json:"duration"`
	LimitQuestion int    `json:"limit_question"`
	Message       BaseResponse
}

type CategoryListResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"value"`
	Bobot int    `json:"bobot"`
}
