package model

type Question struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Value      string `json:"value"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	CreatedBy  string `json:"created_by"`
	UpdatedBy  string `json:"updated_by"`
}
