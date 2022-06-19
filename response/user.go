package response

type UserResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	NIK     string `json:"nik"`
	Message BaseResponse
}

type UserListResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	NIK   string `json:"nik"`
}
