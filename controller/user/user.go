package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/controller"
	"main/controller/auth"
	"main/model"
	"main/response"
	"main/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	rows, err := controller.DB.Query("SELECT id, name, email, nik FROM users ORDER BY id asc")
	if err != nil {
		log.Fatal(err)
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
	id, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}

	var user response.UserListResponse

	err := controller.DB.QueryRow("SELECT id, name, email, nik FROM users WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Email, &user.NIK)
	if err != nil {
		fmt.Print(err)
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

	sqlCekUser := controller.SqlQueryCek("users")
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
		_, err = controller.DB.Exec(sqlStatement, user.Name, user.Email, user.NIK, utils.HashAndSalt([]byte(user.Password)), time.Now(), time.Now())
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
	id, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sqlCekUser := controller.SqlQueryCek("users")
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

		sqlStatement := `UPDATE users SET name = $1, email = $2, nik = $3, password = $4, updated_at = $5 WHERE id = $6`
		_, err = controller.DB.Exec(sqlStatement, user.Name, user.Email, user.NIK, utils.HashAndSalt([]byte(user.Password)), time.Now(), id)
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
	id, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}

	sqlStatement := `DELETE FROM users WHERE id = $1`
	_, err := controller.DB.Exec(sqlStatement, id)
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

func LoginUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	sqlQuery := `SELECT name, email FROM users WHERE email = $1`

	token, err := auth.SignIn(user.Email, user.Password, sqlQuery)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		response.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	response.JSON(w, http.StatusOK, token)
}
