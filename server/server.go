package server

import (
	"main/controller/admin"
	"main/controller/user"

	"github.com/gorilla/mux"
)

func InitializeRoute(r *mux.Router) {

	r.HandleFunc("/login/admin", admin.LoginAdmin).Methods("POST")
	r.HandleFunc("/login/user", user.LoginUser).Methods("POST")

	r.HandleFunc("/admin", admin.GetAdmins).Methods("GET")
	r.HandleFunc("/admin/{id}", admin.GetAdmin).Methods("GET")
	r.HandleFunc("/admin", admin.AddAdmin).Methods("POST")
	r.HandleFunc("/admin/{id}", admin.UpdateAdmin).Methods("POST")
	r.HandleFunc("/admin/{id}", admin.DeleteAdmin).Methods("DELETE")

	r.HandleFunc("/user", user.GetUsers).Methods("GET")
	r.HandleFunc("/user/{id}", user.GetUser).Methods("GET")
	r.HandleFunc("/user", user.AddUser).Methods("POST")
	r.HandleFunc("/user/{id}", user.UpdateUser).Methods("POST")
	r.HandleFunc("/user/{id}", user.DeleteUser).Methods("DELETE")
}
