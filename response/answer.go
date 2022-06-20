package response

type AnswerResponse struct {
	QuestionID int    `json:"question_id"`
	Value      string `json:"value"`
	Score      int    `json:"score"`
	Message    BaseResponse
}

type AnswerListResponse struct {
	Id         string `json:"id"`
	QuestionID int    `json:"question_id"`
	Value      string `json:"value"`
	Score      bool   `json:"score"`
}
