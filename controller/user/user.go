package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var user response.UserResponse

	err = utils.DB.QueryRow("SELECT name, email, nik FROM users WHERE id = $1", id).Scan(&user.Name, &user.Email, &user.NIK)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
		return
	}
	user = response.UserResponse{
		Name:  user.Name,
		Email: user.Email,
		NIK:   user.NIK,
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Get User Successfully!",
		},
	}

	data, _ := json.MarshalIndent(user, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Usr {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var user model.User
	var res response.UserResponse
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sqlCekUser := query.SqlQueryCek(utils.Usr)
	exist := controller.CekUser(user.Email, sqlCekUser)

	if exist.Email != "" && exist.Email != user.Email {

		resp := response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Email Telah Digunakan",
		}
		data, _ := json.MarshalIndent(resp, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return

	}

	sqlGetCurrentPassword := query.SqlGetCurrentPassword(utils.Usr)
	currentPassword := controller.SqlGetCurrentPassword(sqlGetCurrentPassword, id)

	if user.Password == "" {
		user.Password = currentPassword
	} else {
		user.Password = utils.HashAndSalt([]byte(user.Password))
	}

	sqlStatement := `UPDATE users SET name = $1, email = $2, nik = $3, password = $4, updated_at = $5 WHERE id = $6`
	_, err = utils.DB.Exec(sqlStatement, user.Name, user.Email, user.NIK, user.Password, time.Now(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	res = response.UserResponse{
		Name:  user.Name,
		Email: user.Email,
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "User Updated!",
		},
	}
	peopleBytes, _ := json.MarshalIndent(res, "", "\t")
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

	sqlQuery := query.SqlQueryCek(utils.Usr)
	sqlGetID := query.SqlGetID(utils.Usr)
	id := controller.GetUserID(user.Email, sqlGetID)

	name := utils.GetName(int(id), utils.Usr)
	token, err := controller.SignIn(user.Email, user.Password, sqlQuery, utils.Usr, name, id)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		response.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response.JSON(w, http.StatusOK, token)
}
