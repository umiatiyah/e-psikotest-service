package response

type HistoryResponse struct {
	CategoryValue string `json:"category_value"`
	QuestionValue string `json:"question_value"`
	AnswerValue   string `json:"answer_value"`
	AnswerScore   int    `json:"answer_score"`
	User          string `json:"user"`
	NIKUser       string `json:"nik_user"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ValuationResponse struct {
	User          string `json:"user"`
	NIKUser       string `json:"nik_user"`
	CategoryValue string `json:"category_value"`
	TotalScore    int    `json:"total_score"`
}

type ResultResponse struct {
	UserID     int    `json:"user_id"`
	User       string `json:"user"`
	NIKUser    string `json:"nik_user"`
	TotalScore int    `json:"total_score"`
	Result     string `json:"result"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
