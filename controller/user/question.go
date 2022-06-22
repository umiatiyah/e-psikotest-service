package user

import (
	"encoding/json"
	"main/response"
	"main/utils"
	"net/http"
)

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	ps := make(map[int]*response.QuestionUserResponses)

	sql := "SELECT q.id, c.value, q.value FROM question q join category c on c.id = q.category_id WHERE q.is_active = true"
	rows, _ := utils.DB.Query(sql)
	for rows.Next() {
		b := &response.QuestionUserResponses{}
		rows.Scan(&b.ID, &b.Category, &b.Question)

		ps[b.ID] = b
	}

	sql = "SELECT a.id, a.question_id, a.value FROM answer a join question q on q.id = a.question_id WHERE q.is_active = true"
	rows, _ = utils.DB.Query(sql)
	for rows.Next() {
		b := &response.AnswerUserResponses{}
		rows.Scan(&b.ID, &b.QuestionID, &b.Value)

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
