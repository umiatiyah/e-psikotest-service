package response

type QuestionResponse struct {
	CategoryID int    `json:"category_id"`
	Value      string `json:"value"`
	IsActive   bool   `json:"is_active"`
	Message    BaseResponse
}

type QuestionUserResponses struct {
	ID       int                   `json:"id"`
	Category string                `json:"category_name"`
	Question string                `json:"question_value"`
	Answers  []AnswerUserResponses `json:"answers_list"`
}

type AnswerUserResponses struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Value      string `json:"answer_value"`
}
