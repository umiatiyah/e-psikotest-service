package user

import (
	"encoding/json"
	"log"
	"main/controller"
	"main/model"
	"main/response"
	"main/utils"
	"net/http"
	"time"
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

func AddValuation(w http.ResponseWriter, r *http.Request) {
	var history model.DataV2
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var userID int
	var result bool
	var results []bool
	for _, v := range history.Data {

		questionID := controller.GetQuestionIDFromAnswer(v.AnswerID)

		categoryID := controller.GetCategoryIDFromQuestion(questionID)

		sqlStatementHistory := `INSERT INTO history (category_id, question_id, answer_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = utils.DB.Exec(sqlStatementHistory, categoryID, questionID, v.AnswerID, v.UserID, time.Now(), time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		var s model.PreValuation
		err = utils.DB.QueryRow("SELECT c.id, c.min_score, a.score FROM answer a JOIN question q ON a.question_id = q.id JOIN category c ON q.category_id = c.id WHERE a.id = $1", v.AnswerID).
			Scan(&s.CategoryID, &s.MinScore, &s.Score)
		if err != nil {
			log.Fatal(err)
		}

		sqlStatementPrevaluation := `INSERT INTO prevaluation (user_id, category_id, score, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
		_, err = utils.DB.Exec(sqlStatementPrevaluation, v.UserID, s.CategoryID, s.Score, time.Now(), time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		userID = v.UserID

		if s.Score >= s.MinScore {
			results = append(results, true)
		}
		if s.Score < s.MinScore {
			results = append(results, false)
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
	var msg string
	msg = "LULUS"
	if !result {
		msg = "TIDAK LULUS"
	}
	data := response.AnswerResponse{
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: msg,
		},
	}
	dataAnswer, _ := json.MarshalIndent(data, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataAnswer)
}
