package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func GetAdmins(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT id, name, email FROM admin ORDER BY id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var userAdmin []response.AdminListResponse

	for rows.Next() {
		var admin response.AdminListResponse
		rows.Scan(&admin.Id, &admin.Name, &admin.Email)

		userAdmin = append(userAdmin, admin)
	}

	peopleBytes, _ := json.MarshalIndent(userAdmin, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
}

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil || tokenID != uint32(id) || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var admin response.AdminListResponse

	err = utils.DB.QueryRow("SELECT id, name, email FROM admin WHERE id = $1", id).Scan(&admin.Id, &admin.Name, &admin.Email)
	if err != nil {
		fmt.Print(err)
	}

	peopleBytes, _ := json.MarshalIndent(admin, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}

func AddAdmin(w http.ResponseWriter, r *http.Request) {
	var admin model.Admin
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	sqlCekUser := query.SqlQueryCek(utils.Adm)
	exist := controller.CekUser(admin.Email, sqlCekUser)

	if exist.Email != "" {

		userAdmin := response.AdminResponse{
			Name:  admin.Name,
			Email: admin.Email,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Email Telah Digunakan oleh " + exist.Name,
			},
		}
		peopleBytes, _ := json.MarshalIndent(userAdmin, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)

	} else {

		sqlStatement := `INSERT INTO admin (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
		_, err = utils.DB.Exec(sqlStatement, admin.Name, admin.Email, utils.HashAndSalt([]byte(admin.Password)), time.Now(), time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		userAdmin := response.AdminResponse{
			Name:  admin.Name,
			Email: admin.Email,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Admin Created!",
			},
		}
		peopleBytes, _ := json.MarshalIndent(userAdmin, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)

	}
}

func UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil || tokenID != uint32(id) || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var admin model.Admin
	err = json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sqlCekUser := query.SqlQueryCek(utils.Adm)
	exist := controller.CekUser(admin.Email, sqlCekUser)

	if exist.Email != "" && exist.Email != admin.Email {

		userAdmin := response.AdminResponse{
			Name:  admin.Name,
			Email: admin.Email,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Email Telah Digunakan oleh " + exist.Name,
			},
		}
		peopleBytes, _ := json.MarshalIndent(userAdmin, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)
		return

	}
	sqlGetCurrentPassword := query.SqlGetCurrentPassword(utils.Adm)
	currentPassword := controller.SqlGetCurrentPassword(sqlGetCurrentPassword, id)

	if admin.Password == "" {
		admin.Password = currentPassword
	} else {
		admin.Password = utils.HashAndSalt([]byte(admin.Password))
	}

	sqlStatement := `UPDATE admin SET name = $1, email = $2, password = $3, updated_at = $4 WHERE id = $5`
	_, err = utils.DB.Exec(sqlStatement, admin.Name, admin.Email, admin.Password, time.Now(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	userAdmin := response.AdminResponse{
		Name:  admin.Name,
		Email: admin.Email,
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Admin Updated!",
		},
	}
	peopleBytes, _ := json.MarshalIndent(userAdmin, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}

func DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil || tokenID != uint32(id) || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	sqlStatement := `DELETE FROM admin WHERE id = $1`
	_, err = utils.DB.Exec(sqlStatement, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	userAdmin := response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Admin Deleted!",
	}
	peopleBytes, _ := json.MarshalIndent(userAdmin, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	admin := model.Admin{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &admin)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	sqlQuery := query.SqlQueryCek(utils.Adm)
	sqlGetID := query.SqlGetID(utils.Adm)
	id := controller.GetUserID(admin.Email, sqlGetID)

	name := utils.GetName(int(id), utils.Adm)
	token, err := controller.SignIn(admin.Email, admin.Password, sqlQuery, utils.Adm, name, id)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		response.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	if token.Token == "" {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response.JSON(w, http.StatusOK, token)
}
