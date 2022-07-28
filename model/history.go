package model

type DataV2 struct {
	Data []HistoryV2 `json:"data"`
}

type HistoryV2 struct {
	AnswerID  int    `json:"answer_id"`
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Data struct {
	Data []History `json:"data"`
}

type History struct {
	CategoryID int    `json:"category_id"`
	QuestionID int    `json:"question_id"`
	AnswerID   int    `json:"answer_id"`
	UserID     int    `json:"user_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
