package response

type HistoryResponse struct {
	CategoryValue string `json:"category_value"`
	QuestionValue string `json:"question_value"`
	AnswerValue   string `json:"answer_value"`
	User          string `json:"user"`
	NIKUser       string `json:"nik_user"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ValuationResponse struct {
	User          string `json:"user"`
	NIKUser       string `json:"nik_user"`
	CategoryValue string `json:"category_value"`
	MinScore      string `json:"min_score"`
	TotalScore    string `json:"total_score"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ResultResponse struct {
	User      string `json:"user"`
	NIKUser   string `json:"nik_user"`
	Result    string `json:"result"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
