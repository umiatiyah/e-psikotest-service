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

func GetAnswers(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT * FROM answer ORDER BY id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var answers []model.Answer

	for rows.Next() {
		var answer model.Answer
		rows.Scan(&answer.ID, &answer.QuestionID, &answer.Value, &answer.Score, &answer.CreatedAt, &answer.UpdatedAt, &answer.CreatedBy, &answer.UpdatedBy)

		answers = append(answers, answer)
	}

	data, _ := json.MarshalIndent(answers, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	defer rows.Close()
}

func GetAnswer(w http.ResponseWriter, r *http.Request) {
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

	cekIdAnswer := controller.GetMaterialID(int(id), utils.Anw)
	if cekIdAnswer == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Answer Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	var answer model.Answer

	err = utils.DB.QueryRow("SELECT * FROM answer WHERE id = $1", id).
		Scan(&answer.ID, &answer.QuestionID, &answer.Value, &answer.Score,
			&answer.CreatedAt, &answer.CreatedBy, &answer.UpdatedAt, &answer.UpdatedBy)

	if err != nil {
		fmt.Print(err)
	}

	data, _ := json.MarshalIndent(answer, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func AddAnswer(w http.ResponseWriter, r *http.Request) {
	var answer model.Answer
	err := json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	adminName := utils.GetAdminName(int(tokenID), utils.Adm)

	questionID := controller.GetMaterialID(answer.QuestionID, utils.Qst)

	if questionID != answer.QuestionID {
		w.WriteHeader(http.StatusBadRequest)
		res := response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Question Not Found!",
		}
		resError, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(resError)

	} else {

		sqlStatement := `INSERT INTO answer (question_id, value, score, created_at, updated_at, created_by, updated_by ) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = utils.DB.Exec(sqlStatement, answer.QuestionID, answer.Value, answer.Score, time.Now(), time.Now(), adminName, adminName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		data := response.AnswerResponse{
			QuestionID: answer.QuestionID,
			Value:      answer.Value,
			Score:      answer.Score,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Answer Created!",
			},
		}
		dataAnswer, _ := json.MarshalIndent(data, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(dataAnswer)
	}
}

func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
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

	var answer model.Answer
	err = json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	cekIdAnswer := controller.GetMaterialID(int(id), utils.Anw)
	if cekIdAnswer == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Answer Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	adminName := utils.GetAdminName(int(tokenID), utils.Adm)

	questionID := controller.GetMaterialID(answer.QuestionID, utils.Qst)

	if questionID != answer.QuestionID {
		w.WriteHeader(http.StatusBadRequest)
		res := response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Question Not Found!",
		}
		resError, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(resError)

	} else {

		sqlStatement := `UPDATE answer SET question_id = $1, value = $2, score = $3, updated_at = $4, updated_by = $5 WHERE id = $6`
		_, err = utils.DB.Exec(sqlStatement, answer.QuestionID, answer.Value, answer.Score, time.Now(), adminName, id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}
		data := response.AnswerResponse{
			QuestionID: answer.QuestionID,
			Value:      answer.Value,
			Score:      answer.Score,
			Message: response.BaseResponse{
				Status:  http.StatusOK,
				Message: "Answer Updated!",
			},
		}
		dataAnswer, _ := json.MarshalIndent(data, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(dataAnswer)
	}
}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	var answer model.Answer
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

	cekIdAnswer := controller.GetMaterialID(int(id), utils.Anw)
	if cekIdAnswer == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Answer Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	sqlStatement := `DELETE FROM answer WHERE id = $1`
	_, err = utils.DB.Exec(sqlStatement, id)
	row := utils.DB.QueryRow(sqlStatement, id)
	switch err := row.Scan(&answer.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	res := response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Answer Deleted!",
	}
	data, _ := json.MarshalIndent(res, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
