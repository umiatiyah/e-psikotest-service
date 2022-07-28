package response

type AnswerResponse struct {
	ID              int    `json:"id"`
	CategoryID      int    `json:"category_id"`
	CategoryValue   string `json:"category_value"`
	QuestionID      int    `json:"question_id"`
	QuestionValue   string `json:"question_value"`
	QuestiionStatus string `json:"question_status"`
	AnswerValue     string `json:"answer_value"`
	Score           int    `json:"score"`
	Message         BaseResponse
}

type AnswerListResponse struct {
	Id         string `json:"id"`
	QuestionID int    `json:"question_id"`
	Value      string `json:"value"`
	Score      bool   `json:"score"`
}
