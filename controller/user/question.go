package user

import (
	"encoding/json"
	"log"
	"main/controller"
	"main/model"
	"main/response"
	"main/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	ps := make(map[int]*response.QuestionUserResponses)

	sql := "SELECT c.id, c.value, q.id, q.value FROM question q join category c on c.id = q.category_id WHERE q.is_active = true"
	rows, _ := utils.DB.Query(sql)
	for rows.Next() {
		b := &response.QuestionUserResponses{}
		rows.Scan(&b.CategoryID, &b.CategoryValue, &b.QuestionID, &b.QuestionValue)

		ps[b.QuestionID] = b
	}

	sql = "SELECT a.id, a.value, a.question_id FROM answer a join question q on q.id = a.question_id WHERE q.is_active = true"
	rows, _ = utils.DB.Query(sql)
	for rows.Next() {
		b := &response.AnswerUserResponses{}
		rows.Scan(&b.AnswerID, &b.AnswerValue, &b.QuestionID)

		ps[b.QuestionID].Answers = append(ps[b.QuestionID].Answers, *b)
	}

	questions := make([]*response.QuestionUserResponses, 0, len(ps))
	for _, p := range ps {
		questions = append(questions, p)
	}

	data, _ := json.MarshalIndent(questions, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

type User struct {
	ID int `json:"id"`
}

func CekUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var IDs []int
	sql := "SELECT user_id FROM valuation WHERE user_id = $1"
	rows, err := utils.DB.Query(sql, id)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var idUser int
		rows.Scan(&idUser)
		IDs = append(IDs, idUser)
	}
	res := response.BaseResponse{}
	if len(IDs) > 0 {
		res = response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Exists",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}
	res = response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Not Exists",
	}
	data, _ := json.MarshalIndent(res, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func SaveHistory(w http.ResponseWriter, r *http.Request) {
	var history model.Data
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for _, v := range history.Data {
		sqlStatement := `INSERT INTO history (category_id, question_id, answer_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = utils.DB.Exec(sqlStatement, v.CategoryID, v.QuestionID, v.AnswerID, v.UserID, time.Now(), time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}
	}
	data := response.AnswerResponse{
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Answer Created!",
		},
	}
	dataAnswer, _ := json.MarshalIndent(data, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataAnswer)
}

type BO struct {
	Bobot    int
	Category int
}

func AddValuation(w http.ResponseWriter, r *http.Request) {
	var history model.DataV2
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var userID int
	var result bool
	var results []bool
	bos := make(map[int]BO)
	bobotScore := make(map[int]int)
	category := make(map[int]int)
	for _, v := range history.Data {

		questionID := controller.GetQuestionIDFromAnswer(v.AnswerID)

		categoryID := controller.GetCategoryIDFromQuestion(questionID)

		sqlStatementHistory := `INSERT INTO history (category_id, question_id, answer_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = utils.DB.Exec(sqlStatementHistory, categoryID, questionID, v.AnswerID, v.UserID, time.Now(), time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}
		userID = v.UserID
	}

	var bobotUser []model.PreBobotValuation
	bobot, err := utils.DB.Query("SELECT pv.user_id, c.id, a.score FROM history pv JOIN answer a ON pv.answer_id = a.id JOIN question q ON a.question_id = q.id JOIN category c ON q.category_id = c.id JOIN users u ON pv.user_id = u.id WHERE pv.user_id = $1 GROUP BY pv.user_id, c.id, a.score", userID)
	if err != nil {
		log.Fatal(err)
	}

	for bobot.Next() {
		var bv model.PreBobotValuation
		bobot.Scan(&bv.UserID, &bv.CategoryID, &bv.Score)

		bobotUser = append(bobotUser, bv)
	}

	for _, v := range bobotUser {
		bobotScore[v.CategoryID] += v.Score
		category[v.CategoryID] = v.CategoryID
		bos[v.CategoryID] = BO{
			Bobot:    bobotScore[v.CategoryID],
			Category: v.CategoryID,
		}
	}

	for _, v := range bos {
		sqlStatementValuation := `INSERT INTO bobotvaluation (user_id, category_id, total_score, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
		_, err = utils.DB.Exec(sqlStatementValuation, userID, v.Category, v.Bobot, time.Now(), time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		var s model.PreValuation
		err = utils.DB.QueryRow("SELECT c.id, c.min_score FROM category c WHERE c.id = $1", v.Category).
			Scan(&s.CategoryID, &s.MinScore)
		if err != nil {
			log.Fatal(err)
		}
		if s.CategoryID == v.Category {
			if v.Bobot >= s.MinScore {
				results = append(results, true)
			}
			if v.Bobot < s.MinScore {
				results = append(results, false)
			}
		}
	}

	result = true
	for _, v := range results {
		if !v {
			result = false
		}
	}
	sqlStatementValuation := `INSERT INTO valuation (user_id, result, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err = utils.DB.Exec(sqlStatementValuation, userID, result, time.Now(), time.Now())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	data := response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Test Completed!",
	}
	dataAnswer, _ := json.MarshalIndent(data, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataAnswer)
}

func AddValuationNew(w http.ResponseWriter, r *http.Request) {
	var history model.DataV2
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var userID int
	for _, v := range history.Data {

		questionID := controller.GetQuestionIDFromAnswer(v.AnswerID)

		categoryID := controller.GetCategoryIDFromQuestion(questionID)

		sqlStatementHistory := `INSERT INTO history (category_id, question_id, answer_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = utils.DB.Exec(sqlStatementHistory, categoryID, questionID, v.AnswerID, v.UserID, time.Now(), time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}
		userID = v.UserID
	}

	var JumlahItem int
	item := utils.DB.QueryRow("SELECT COUNT(*) FROM question WHERE is_active = true").
		Scan(&JumlahItem)
	if item != nil {
		log.Fatal(item)
	}

	BobotMax := 5
	BobotMin := 1
	Xminimal := JumlahItem * BobotMin
	Xmaximal := JumlahItem * BobotMax
	Range := Xmaximal - Xminimal
	Mean := (Xmaximal + Xminimal) / 2
	SD := Range / (BobotMax + BobotMin)

	var X int
	var res string
	for _, v := range history.Data {

		score := utils.DB.QueryRow("SELECT SUM(a.score) FROM history h JOIN answer a ON h.answer_id = a.id WHERE user_id = $1", v.UserID).
			Scan(&X)
		if score != nil {
			log.Fatal(score)
		}
	}

	s1 := float64(Mean) - (1.5 * float64(SD))
	s2 := float64(Mean) - (0.5 * float64(SD))
	s3 := float64(Mean) + (0.5 * float64(SD))
	s4 := float64(Mean) + (1.5 * float64(SD))
	if float64(X) <= s1 {
		res = "Sangat Rendah"
	} else if s1 < float64(X) && float64(X) <= s2 {
		res = "Rendah"
	} else if s2 < float64(X) && float64(X) <= s3 {
		res = "Sedang"
	} else if s3 < float64(X) && float64(X) <= s4 {
		res = "Tinggi"
	} else if s4 < float64(X) {
		res = "Sangat Tinggi"
	}

	sqlStatementValuation := `INSERT INTO valuation (user_id, result, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err = utils.DB.Exec(sqlStatementValuation, userID, res, time.Now(), time.Now())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	data := response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Test Completed!",
	}
	dataAnswer, _ := json.MarshalIndent(data, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataAnswer)
}
