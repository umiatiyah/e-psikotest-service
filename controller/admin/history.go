package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/controller/auth"
	"main/response"
	"main/utils"
	"net/http"
	"sort"
)

func GetHistory(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT c.value, q.value, a.value, a.score, u.name, u.nik, h.created_at, h.updated_at FROM history h JOIN category c ON h.category_id = c.id JOIN question q ON h.question_id = q.id JOIN answer a ON h.answer_id = a.id JOIN users u ON h.user_id = u.id ORDER BY h.created_at desc")
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
		rows.Scan(&history.CategoryValue, &history.QuestionValue, &history.AnswerValue, &history.AnswerScore, &history.User, &history.NIKUser, &history.CreatedAt, &history.UpdatedAt)

		histories = append(histories, history)
	}

	peopleBytes, _ := json.MarshalIndent(histories, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
}

func GetResult(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT u.id, u.name, u.nik, v.result, v.created_at, v.updated_at FROM valuation v JOIN users u ON v.user_id = u.id ORDER BY v.id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var results []response.ResultResponse

	for rows.Next() {
		var result response.ResultResponse
		rows.Scan(&result.UserID, &result.User, &result.NIKUser, &result.Result, &result.CreatedAt, &result.UpdatedAt)

		score := utils.DB.QueryRow("SELECT SUM(a.score) FROM history h JOIN answer a ON h.answer_id = a.id WHERE user_id = $1", result.UserID).
			Scan(&result.TotalScore)
		if score != nil {
			log.Fatal(score)
		}
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalScore > results[j].TotalScore
	})

	peopleBytes, _ := json.MarshalIndent(results, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
}

func GetValuation(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT u.name, u.nik, c.value, (SUM(a.score)) as bobot FROM history h JOIN category c ON h.category_id = c.id JOIN answer a ON h.answer_id = a.id JOIN users u ON h.user_id = u.id GROUP BY u.name, u.nik, c.value, date(h.created_at) ORDER BY date(h.created_at) desc")
	if err != nil {
		log.Fatal(err)
	}

	totalScore, err := utils.DB.Query("SELECT (SUM(a.score)) as bobot FROM history h JOIN category c ON h.category_id = c.id JOIN answer a ON h.answer_id = a.id JOIN users u ON h.user_id = u.id GROUP BY u.name, u.nik, c.value ORDER BY bobot DESC")
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
	var totalJumlahNilaiPerkategori, jumlahNilaiPerkategori int

	for totalScore.Next() {
		var i int
		totalScore.Scan(&i)
		totalJumlahNilaiPerkategori += i
	}

	for rows.Next() {
		var valuation response.ValuationResponse
		rows.Scan(&valuation.User, &valuation.NIKUser, &valuation.CategoryValue, &valuation.TotalScore)
		jumlahNilaiPerkategori = valuation.TotalScore
		u := (float64(jumlahNilaiPerkategori)) / (float64(totalJumlahNilaiPerkategori)) * 100

		s := fmt.Sprintf("%.2f", u)
		valuation.PersenBobotCategory = s
		valuations = append(valuations, valuation)
	}

	peopleBytes, _ := json.MarshalIndent(valuations, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
}
