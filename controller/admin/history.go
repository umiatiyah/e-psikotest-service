package admin

import (
	"encoding/json"
	"errors"
	"log"
	"main/controller/auth"
	"main/response"
	"main/utils"
	"net/http"
)

func GetHistory(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT c.value, q.value, a.value, u.name, u.nik, h.created_at, h.updated_at FROM history h JOIN category c ON h.category_id = c.id JOIN question q ON h.question_id = q.id JOIN answer a ON h.answer_id = a.id JOIN users u ON h.user_id = u.id ORDER BY h.id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var histories []response.HistoryResponse

	for rows.Next() {
		var history response.HistoryResponse
		rows.Scan(&history.CategoryValue, &history.QuestionValue, &history.AnswerValue, &history.User, &history.NIKUser, &history.CreatedAt, &history.UpdatedAt)

		histories = append(histories, history)
	}

	peopleBytes, _ := json.MarshalIndent(histories, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
}

func GetValuation(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT u.name, u.nik, v.result, v.created_at, v.updated_at FROM valuation v JOIN users u ON v.user_id = u.id ORDER BY v.id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var valuations []response.ValuationResponse

	for rows.Next() {
		var valuation response.ValuationResponse
		rows.Scan(&valuation.User, &valuation.NIKUser, &valuation.Result, &valuation.CreatedAt, &valuation.UpdatedAt)

		if valuation.Result == "true" {
			valuation.Result = "LULUS"
		} else {
			valuation.Result = "TIDAK LULUS"
		}

		valuations = append(valuations, valuation)
	}

	peopleBytes, _ := json.MarshalIndent(valuations, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
}
