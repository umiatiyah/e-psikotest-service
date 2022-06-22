package response

type Dashboard struct {
	CountAdmin    int `json:"count_admin"`
	CountUser     int `json:"count_user"`
	CountCategory int `json:"count_category"`
	CountQuestion int `json:"count_question"`
	CountAnswer   int `json:"count_answer"`
}
