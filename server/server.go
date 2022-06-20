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

	r.HandleFunc("/admin", auth.MiddlewareAuth(admin.GetAdmins)).Methods("GET")
	r.HandleFunc("/admin/{id}", auth.MiddlewareAuth(admin.GetAdmin)).Methods("GET")
	r.HandleFunc("/admin", auth.MiddlewareAuth(admin.AddAdmin)).Methods("POST")
	r.HandleFunc("/admin/{id}", auth.MiddlewareAuth(admin.UpdateAdmin)).Methods("POST")
	r.HandleFunc("/admin/{id}", auth.MiddlewareAuth(admin.DeleteAdmin)).Methods("DELETE")

	r.HandleFunc("/user", auth.MiddlewareAuth(user.GetUsers)).Methods("GET")
	r.HandleFunc("/user/{id}", auth.MiddlewareAuth(user.GetUser)).Methods("GET")
	r.HandleFunc("/user", auth.MiddlewareAuth(user.AddUser)).Methods("POST")
	r.HandleFunc("/user/{id}", auth.MiddlewareAuth(user.UpdateUser)).Methods("POST")
	r.HandleFunc("/user/{id}", auth.MiddlewareAuth(user.DeleteUser)).Methods("DELETE")

	r.HandleFunc("/category", auth.MiddlewareAuth(admin.GetCategories)).Methods("GET")
	r.HandleFunc("/category/{id}", auth.MiddlewareAuth(admin.GetCategory)).Methods("GET")
	r.HandleFunc("/category", auth.MiddlewareAuth(admin.AddCategory)).Methods("POST")
	r.HandleFunc("/category/{id}", auth.MiddlewareAuth(admin.UpdateCategory)).Methods("POST")
	r.HandleFunc("/category/{id}", auth.MiddlewareAuth(admin.DeleteCategory)).Methods("DELETE")

	r.HandleFunc("/question", auth.MiddlewareAuth(admin.GetQuestions)).Methods("GET")
	r.HandleFunc("/question/{id}", auth.MiddlewareAuth(admin.GetQuestion)).Methods("GET")
	r.HandleFunc("/question", auth.MiddlewareAuth(admin.AddQuestion)).Methods("POST")
	r.HandleFunc("/question/{id}", auth.MiddlewareAuth(admin.UpdateQuestion)).Methods("POST")
	r.HandleFunc("/question/{id}", auth.MiddlewareAuth(admin.DeleteQuestion)).Methods("DELETE")
}
