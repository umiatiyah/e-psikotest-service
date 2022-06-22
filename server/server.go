package server

import (
	"main/controller/admin"
	"main/controller/auth"
	"main/controller/user"

	"github.com/gorilla/mux"
)

func InitializeRoute(r *mux.Router) {

	r.HandleFunc("/login/admin", admin.LoginAdmin).Methods("POST")
	r.HandleFunc("/login/user", user.LoginUser).Methods("POST")

	r.HandleFunc("/admin/dashboard", auth.MiddlewareAuth(admin.Dashboard)).Methods("GET")

	r.HandleFunc("/admin/list", auth.MiddlewareAuth(admin.GetAdmins)).Methods("GET")
	r.HandleFunc("/admin/detail/{id}", auth.MiddlewareAuth(admin.GetAdmin)).Methods("GET")
	r.HandleFunc("/admin/create", auth.MiddlewareAuth(admin.AddAdmin)).Methods("POST")
	r.HandleFunc("/admin/update/{id}", auth.MiddlewareAuth(admin.UpdateAdmin)).Methods("POST")
	r.HandleFunc("/admin/delete/{id}", auth.MiddlewareAuth(admin.DeleteAdmin)).Methods("DELETE")

	r.HandleFunc("/admin/user", auth.MiddlewareAuth(user.GetUsers)).Methods("GET")
	r.HandleFunc("/admin/user/{id}", auth.MiddlewareAuth(user.GetUser)).Methods("GET")
	r.HandleFunc("/admin/user", auth.MiddlewareAuth(user.AddUser)).Methods("POST")
	r.HandleFunc("/admin/user/{id}", auth.MiddlewareAuth(user.UpdateUser)).Methods("POST")
	r.HandleFunc("/admin/user/{id}", auth.MiddlewareAuth(user.DeleteUser)).Methods("DELETE")

	r.HandleFunc("/admin/category", auth.MiddlewareAuth(admin.GetCategories)).Methods("GET")
	r.HandleFunc("/admin/category/{id}", auth.MiddlewareAuth(admin.GetCategory)).Methods("GET")
	r.HandleFunc("/admin/category", auth.MiddlewareAuth(admin.AddCategory)).Methods("POST")
	r.HandleFunc("/admin/category/{id}", auth.MiddlewareAuth(admin.UpdateCategory)).Methods("POST")
	r.HandleFunc("/admin/category/{id}", auth.MiddlewareAuth(admin.DeleteCategory)).Methods("DELETE")

	r.HandleFunc("/admin/question", auth.MiddlewareAuth(admin.GetQuestions)).Methods("GET")
	r.HandleFunc("/admin/question/{id}", auth.MiddlewareAuth(admin.GetQuestion)).Methods("GET")
	r.HandleFunc("/admin/question", auth.MiddlewareAuth(admin.AddQuestion)).Methods("POST")
	r.HandleFunc("/admin/question/{id}", auth.MiddlewareAuth(admin.UpdateQuestion)).Methods("POST")
	r.HandleFunc("/admin/question/{id}", auth.MiddlewareAuth(admin.DeleteQuestion)).Methods("DELETE")

	r.HandleFunc("/admin/answer", auth.MiddlewareAuth(admin.GetAnswers)).Methods("GET")
	r.HandleFunc("/admin/answer/{id}", auth.MiddlewareAuth(admin.GetAnswer)).Methods("GET")
	r.HandleFunc("/admin/answer", auth.MiddlewareAuth(admin.AddAnswer)).Methods("POST")
	r.HandleFunc("/admin/answer/{id}", auth.MiddlewareAuth(admin.UpdateAnswer)).Methods("POST")
	r.HandleFunc("/admin/answer/{id}", auth.MiddlewareAuth(admin.DeleteAnswer)).Methods("DELETE")

	r.HandleFunc("/user/question", auth.MiddlewareAuth(user.GetQuestions)).Methods("GET")
}
