package model

type Valuation struct {
	UserID    int    `json:"user_id"`
	Score     int    `json:"total_score"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PreBobotValuation struct {
	UserID     int `json:"user_id"`
	CategoryID int `json:"category_id"`
	Score      int `json:"score"`
}

type BobotValuation struct {
	UserID     int    `json:"user_id"`
	CategoryID int    `json:"category_id"`
	TotalScore int    `json:"total_score"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PreValuation struct {
	CategoryID int `json:"category_id"`
	MinScore   int `json:"min_score"`
}
