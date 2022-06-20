package model

type Answer struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Value      string `json:"value"`
	Score      int    `json:"score"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	CreatedBy  string `json:"created_by"`
	UpdatedBy  string `json:"updated_by"`
}
