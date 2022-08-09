package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/controller"
	"main/controller/auth"
	"main/query"
	"main/response"
	"main/utils"
	"math"
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

func SqlGetMaxBobotCategory(category string) int {
	SqlCreateTempTblBobot := query.SqlCreateTempTblBobot()
	controller.SqlCreateTempTblBobot(SqlCreateTempTblBobot)
	SqlGetMaxBobotCategory := query.SqlGetMaxBobotCategory()
	return controller.SqlGetMaxBobotCategory(SqlGetMaxBobotCategory, category)
}

func GetResult(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT u.id, u.name, u.nik, v.result, v.created_at, v.updated_at FROM valuation v JOIN users u ON v.user_id = u.id ORDER BY v.id asc")
	if err != nil {
		log.Fatal(err)
	}

	maxBobotDataDiri := SqlGetMaxBobotCategory(utils.CDataDiri)
	maxBobotPengalaman := SqlGetMaxBobotCategory(utils.CPengalaman)
	maxBobotLoyalitas := SqlGetMaxBobotCategory(utils.CLoyalitas)
	maxBobotTekananKerja := SqlGetMaxBobotCategory(utils.CTekananKerja)
	maxBobotMotivasi := SqlGetMaxBobotCategory(utils.CMotivasi)
	maxBobotSkill := SqlGetMaxBobotCategory(utils.CSkill)
	maxBobotPotensi := SqlGetMaxBobotCategory(utils.CPotensi)

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	type Test struct {
		User int
		Category string
		Bobot int
	}
	var results []response.ResultResponse
	var tests []Test

	for rows.Next() {
		var result response.ResultResponse
		rows.Scan(&result.UserID, &result.User, &result.NIKUser, &result.Result, &result.CreatedAt, &result.UpdatedAt)

		bobotUserCategory, err := utils.DB.Query("select * from tempBobot where userid = $1 ", result.UserID)
		if err != nil {
			log.Fatal(err)
		}

		for bobotUserCategory.Next() {
			var test Test
			bobotUserCategory.Scan(&test.User,&test.Category, &test.Bobot)
			tests = append(tests, test)
		}

		var c1,c2,c3,c4,c5,c6,c7 float64
		for _, v := range tests {
			if v.Category == utils.CDataDiri {
				c1 = float64(v.Bobot) / float64(maxBobotDataDiri)
				c1 = (math.Round(c1*100)/100)*12
			}
			if v.Category == utils.CPengalaman {
				c2 = float64(v.Bobot) / float64(maxBobotPengalaman)
				c2 = (math.Round(c2*100)/100)*20
			}
			if v.Category == utils.CLoyalitas {
				c3 = float64(v.Bobot) / float64(maxBobotLoyalitas)
				c3 = (math.Round(c3*100)/100)*12
			}
			if v.Category == utils.CTekananKerja {
				c4 = float64(v.Bobot) / float64(maxBobotTekananKerja)
				c4 = (math.Round(c4*100)/100)*16
			}
			if v.Category == utils.CMotivasi {
				c5 = float64(v.Bobot) / float64(maxBobotMotivasi)
				c5 = (math.Round(c5*100)/100)*8
			}
			if v.Category == utils.CSkill {
				c6 = float64(v.Bobot) / float64(maxBobotSkill)
				c6 = (math.Round(c6*100)/100)*8
			}
			if v.Category == utils.CPotensi {
				c7 = float64(v.Bobot) / float64(maxBobotPotensi)
				c7 = (math.Round(c7*100)/100)*8
			}
		total := c1+c2+c3+c4+c5+c6+c7
		s := fmt.Sprintf("%.2f", total)
		result.TotalScore = s
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
