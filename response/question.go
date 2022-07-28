package response

type QuestionResponse struct {
	ID           int    `json:"id"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	Value        string `json:"value"`
	IsActive     string `json:"is_active"`
	Message      BaseResponse
}

type QuestionUserResponses struct {
	CategoryID    int                   `json:"category_id"`
	CategoryValue string                `json:"category_value"`
	QuestionID    int                   `json:"question_id"`
	QuestionValue string                `json:"question_value"`
	Answers       []AnswerUserResponses `json:"answers_list"`
}

type AnswerUserResponses struct {
	AnswerID    int    `json:"answer_id"`
	AnswerValue string `json:"answer_value"`
	QuestionID  int    `json:"question_id"`
}
