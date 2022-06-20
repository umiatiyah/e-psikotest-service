package model

type Category struct {
	ID            int    `json:"id"`
	Value         string `json:"value"`
	MinScore      int    `json:"min_score"`
	Duration      int    `json:"duration"`
	LimitQuestion int    `json:"limit_question"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	CreatedBy     string `json:"created_by"`
	UpdatedBy     string `json:"updated_by"`
}
