package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/controller"
	"main/controller/auth"
	"main/model"
	"main/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {

	rows, err := controller.DB.Query("SELECT * FROM category ORDER BY id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if role != "admin" {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var categories []model.Category

	for rows.Next() {
		var category model.Category
		rows.Scan(&category.ID, &category.Value, &category.MinScore, &category.Duration, &category.LimitQuestion,
			&category.CreatedAt, &category.UpdatedAt, &category.CreatedBy, &category.UpdatedBy)

		categories = append(categories, category)
	}

	peopleBytes, _ := json.MarshalIndent(categories, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if role != "admin" {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var category model.Category

	err = controller.DB.QueryRow("SELECT * FROM category WHERE id = $1", id).
		Scan(&category.ID, &category.Value, &category.MinScore, &category.Duration, &category.LimitQuestion,
			&category.CreatedAt, &category.CreatedBy, &category.UpdatedAt, &category.UpdatedBy)

	if err != nil {
		fmt.Print(err)
	}

	peopleBytes, _ := json.MarshalIndent(category, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if role != "admin" {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	adminName := controller.GetAdminName(int(tokenID))

	sqlStatement := `INSERT INTO category (value, min_score, duration, limit_question, created_at, updated_at, created_by, updated_by ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = controller.DB.Exec(sqlStatement, category.Value, category.MinScore, category.Duration, category.LimitQuestion, time.Now(), time.Now(), adminName, adminName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	data := response.CategoryResponse{
		Value:         category.Value,
		MinScore:      category.MinScore,
		Duration:      category.Duration,
		LimitQuestion: category.LimitQuestion,
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Category Created!",
		},
	}
	peopleBytes, _ := json.MarshalIndent(data, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(id) && role != "admin" {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var category model.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	adminName := controller.GetAdminName(int(tokenID))

	sqlStatement := `UPDATE category SET value = $1, min_score = $2, duration = $3, limit_question = $4, updated_at = $5, updated_by = $6 WHERE id = $7`
	_, err = controller.DB.Exec(sqlStatement, category.Value, category.MinScore, category.Duration, category.LimitQuestion, time.Now(), adminName, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	data := response.CategoryResponse{
		Value:         category.Value,
		MinScore:      category.MinScore,
		Duration:      category.Duration,
		LimitQuestion: category.LimitQuestion,
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Category Updated!",
		},
	}
	peopleBytes, _ := json.MarshalIndent(data, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, role, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(id) && role != "admin" {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	sqlStatement := `DELETE FROM category WHERE id = $1`
	_, err = controller.DB.Exec(sqlStatement, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	res := response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Category Deleted!",
	}
	peopleBytes, _ := json.MarshalIndent(res, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

}