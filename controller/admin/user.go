package admin

import (
	"encoding/json"
	"errors"
	"log"
	"main/controller"
	"main/controller/auth"
	"main/model"
	"main/query"
	"main/response"
	"main/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT id, name, email, nik FROM users ORDER BY id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var users []response.UserListResponse

	for rows.Next() {
		var user response.UserListResponse
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.NIK)

		users = append(users, user)
	}

	peopleBytes, _ := json.MarshalIndent(users, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var user response.UserListResponse

	err = utils.DB.QueryRow("SELECT id, name, email, nik FROM users WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Email, &user.NIK)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}

	peopleBytes, _ := json.MarshalIndent(user, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	sqlCekUser := query.SqlQueryCek(utils.Usr)
	exist := controller.CekUser(user.Email, sqlCekUser)

	if exist.Email != "" {

		user := response.UserResponse{
			Name:  user.Name,
			Email: user.Email,
			NIK:   user.NIK,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Email Telah Digunakan oleh " + exist.Name,
			},
		}
		peopleBytes, _ := json.MarshalIndent(user, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)

	} else {

		sqlStatement := `INSERT INTO users (name, email, nik, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = utils.DB.Exec(sqlStatement, user.Name, user.Email, user.NIK, utils.HashAndSalt([]byte(user.Password)), time.Now(), time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		user := response.UserResponse{
			Name:  user.Name,
			Email: user.Email,
			NIK:   user.NIK,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "User Created!",
			},
		}
		peopleBytes, _ := json.MarshalIndent(user, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)

	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var user model.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sqlCekUser := query.SqlQueryCek(utils.Usr)
	exist := controller.CekUser(user.Email, sqlCekUser)

	if exist.Email != "" && exist.Email != user.Email {

		user := response.UserResponse{
			Name:  user.Name,
			Email: user.Email,
			NIK:   user.NIK,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Email Telah Digunakan oleh " + exist.Name,
			},
		}
		peopleBytes, _ := json.MarshalIndent(user, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)

	} else {

		sqlStatement := `UPDATE users SET name = $1, email = $2, nik = $3, password = $4, updated_at = $5 WHERE id = $6`
		_, err = utils.DB.Exec(sqlStatement, user.Name, user.Email, user.NIK, utils.HashAndSalt([]byte(user.Password)), time.Now(), id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}
		user := response.UserResponse{
			Name:  user.Name,
			Email: user.Email,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "User Updated!",
			},
		}
		peopleBytes, _ := json.MarshalIndent(user, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)

	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	sqlStatement := `DELETE FROM users WHERE id = $1`
	_, err = utils.DB.Exec(sqlStatement, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	user := response.BaseResponse{
		Status:  http.StatusOK,
		Message: "User Deleted!",
	}
	peopleBytes, _ := json.MarshalIndent(user, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}
