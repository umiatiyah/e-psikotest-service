package response

type QuestionResponse struct {
	CategoryID int    `json:"category_id"`
	Value      string `json:"value"`
	IsActive   bool   `json:"is_active"`
	Message    BaseResponse
}

type QuestionListResponse struct {
	Id         string `json:"id"`
	CategoryID int    `json:"category_id"`
	Value      string `json:"value"`
	IsActive   bool   `json:"is_active"`
}
