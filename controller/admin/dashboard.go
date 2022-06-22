package admin

import (
	"encoding/json"
	"errors"
	"log"
	"main/controller/auth"
	"main/query"
	"main/response"
	"main/utils"
	"net/http"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var count response.Dashboard
	err = utils.DB.QueryRow(query.SqlCount(utils.Adm)).Scan(&count.CountAdmin)
	if err != nil {
		log.Print(err)
	}
	err = utils.DB.QueryRow(query.SqlCount(utils.Usr)).Scan(&count.CountUser)
	if err != nil {
		log.Print(err)
	}
	err = utils.DB.QueryRow(query.SqlCount(utils.Ctg)).Scan(&count.CountCategory)
	if err != nil {
		log.Print(err)
	}
	err = utils.DB.QueryRow(query.SqlCount(utils.Qst)).Scan(&count.CountQuestion)
	if err != nil {
		log.Print(err)
	}
	err = utils.DB.QueryRow(query.SqlCount(utils.Anw)).Scan(&count.CountAnswer)
	if err != nil {
		log.Print(err)
	}

	peopleBytes, _ := json.MarshalIndent(count, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}
