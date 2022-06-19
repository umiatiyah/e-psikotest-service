package response

type AdminResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message BaseResponse
}

type AdminListResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
