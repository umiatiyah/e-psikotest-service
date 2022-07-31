package admin

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/controller"
	"main/controller/auth"
	"main/model"
	"main/response"
	"main/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetQuestions(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT q.id, q.category_id, c.value, q.value, q.is_active FROM question q JOIN category c ON q.category_id = c.id ORDER BY q.id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var questions []response.QuestionResponse

	for rows.Next() {
		var question response.QuestionResponse
		rows.Scan(&question.ID, &question.CategoryID, &question.CategoryName, &question.Value, &question.IsActive)

		if question.IsActive == "true" {
			question.IsActive = "Active"
		} else {
			question.IsActive = "Inactive"
		}

		questions = append(questions, question)
	}

	data, _ := json.MarshalIndent(questions, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	defer rows.Close()
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
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

	cekIdQuestion := controller.GetMaterialID(int(id), utils.Qst)
	if cekIdQuestion == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Question Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	var question model.Question

	err = utils.DB.QueryRow("SELECT q.id, q.category_id, c.value, q.value, q.is_active, q.created_at, q.created_by, q.updated_at, q.updated_by FROM question q JOIN category c ON q.category_id = c.id WHERE q.id = $1", id).
		Scan(&question.ID, &question.CategoryID, &question.CategoryName, &question.Value, &question.IsActive,
			&question.CreatedAt, &question.CreatedBy, &question.UpdatedAt, &question.UpdatedBy)

	if err != nil {
		fmt.Print(err)
	}

	data, _ := json.MarshalIndent(question, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func AddQuestion(w http.ResponseWriter, r *http.Request) {
	var question model.Question
	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	adminName := utils.GetName(int(tokenID), utils.Adm)

	categoryID := controller.GetMaterialID(question.CategoryID, utils.Ctg)

	if categoryID != question.CategoryID {
		w.WriteHeader(http.StatusBadRequest)
		res := response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Category Not Found!",
		}
		resError, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(resError)

	} else {

		sqlStatement := `INSERT INTO question (category_id, value, is_active, created_at, updated_at, created_by, updated_by ) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = utils.DB.Exec(sqlStatement, question.CategoryID, question.Value, question.IsActive, time.Now(), time.Now(), adminName, adminName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		data := response.QuestionResponse{
			CategoryID: question.CategoryID,
			Value:      question.Value,
			IsActive:   strconv.FormatBool(question.IsActive),
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Questiion Created!",
			},
		}
		dataQuestion, _ := json.MarshalIndent(data, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(dataQuestion)
	}
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var question model.Question
	err = json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	cekIdQuestion := controller.GetMaterialID(int(id), utils.Qst)
	if cekIdQuestion == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Question Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	adminName := utils.GetName(int(tokenID), utils.Adm)

	categoryID := controller.GetMaterialID(question.CategoryID, utils.Ctg)

	if categoryID != question.CategoryID {
		w.WriteHeader(http.StatusBadRequest)
		res := response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Category Not Found!",
		}
		resError, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(resError)

	} else {

		sqlStatement := `UPDATE question SET category_id = $1, value = $2, is_active = $3, updated_at = $4, updated_by = $5 WHERE id = $6`
		_, err = utils.DB.Exec(sqlStatement, question.CategoryID, question.Value, question.IsActive, time.Now(), adminName, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}
		data := response.QuestionResponse{
			CategoryID: question.CategoryID,
			Value:      question.Value,
			IsActive:   strconv.FormatBool(question.IsActive),
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Question Updated!",
			},
		}
		dataQuestion, _ := json.MarshalIndent(data, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(dataQuestion)
	}
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	var question model.Question
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

	cekIdQuestion := controller.GetMaterialID(int(id), utils.Qst)
	if cekIdQuestion == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Question Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	CekQuestionInAnswer := controller.CekMaterialInOtherRelation(int(id), "question_id", utils.Anw)
	if CekQuestionInAnswer == int(id) {
		res := response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Question Used In Answer!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}
	sqlStatement := `DELETE FROM question WHERE id = $1`
	_, err = utils.DB.Exec(sqlStatement, id)
	row := utils.DB.QueryRow(sqlStatement, id)
	switch err := row.Scan(&question.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	res := response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Question Deleted!",
	}
	data, _ := json.MarshalIndent(res, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
