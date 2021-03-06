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

func GetCategories(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("SELECT * FROM category ORDER BY id asc")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var categories []model.Category

	for rows.Next() {
		var category model.Category
		rows.Scan(&category.ID, &category.Name, &category.MinScore, &category.Duration, &category.LimitQuestion,
			&category.CreatedAt, &category.UpdatedAt, &category.CreatedBy, &category.UpdatedBy)

		categories = append(categories, category)
	}

	data, _ := json.MarshalIndent(categories, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	defer rows.Close()
}

func GetBobotCategories(w http.ResponseWriter, r *http.Request) {

	rows, err := utils.DB.Query("select c.id, c.value, count(a.*) as bobot FROM answer a JOIN question q ON a.question_id = q.id JOIN category c ON q.category_id = c.id group by c.id order by c.id")
	if err != nil {
		log.Fatal(err)
	}

	_, role, err := auth.ExtractTokenID(r)
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	var categories []response.CategoryListResponse

	for rows.Next() {
		var category response.CategoryListResponse
		rows.Scan(&category.ID, &category.Name, &category.Bobot)
		categories = append(categories, category)
	}

	data, _ := json.MarshalIndent(categories, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

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
	if err != nil || role != utils.Adm {
		w.Header().Set("Content-Type", "application/json")
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	cekIdCategory := controller.GetMaterialID(int(id), utils.Ctg)
	if cekIdCategory == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Category Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	var category model.Category

	err = utils.DB.QueryRow("SELECT * FROM category WHERE id = $1", id).
		Scan(&category.ID, &category.Name, &category.MinScore, &category.Duration, &category.LimitQuestion,
			&category.CreatedAt, &category.CreatedBy, &category.UpdatedAt, &category.UpdatedBy)

	if err != nil {
		fmt.Print(err)
	}

	data, _ := json.MarshalIndent(category, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
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

	sqlStatement := `INSERT INTO category (value, min_score, duration, limit_question, created_at, updated_at, created_by, updated_by ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = utils.DB.Exec(sqlStatement, category.Name, category.MinScore, category.Duration, category.LimitQuestion, time.Now(), time.Now(), adminName, adminName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	data := response.CategoryResponse{
		Name:          category.Name,
		MinScore:      category.MinScore,
		Duration:      category.Duration,
		LimitQuestion: category.LimitQuestion,
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Category Created!",
		},
	}
	dataCategory, _ := json.MarshalIndent(data, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataCategory)

}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
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

	var category model.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	cekIdCategory := controller.GetMaterialID(int(id), utils.Ctg)
	if cekIdCategory == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Category Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	adminName := utils.GetName(int(tokenID), utils.Adm)

	sqlStatement := `UPDATE category SET value = $1, min_score = $2, duration = $3, limit_question = $4, updated_at = $5, updated_by = $6 WHERE id = $7`
	_, err = utils.DB.Exec(sqlStatement, category.Name, category.MinScore, category.Duration, category.LimitQuestion, time.Now(), adminName, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	data := response.CategoryResponse{
		Name:          category.Name,
		MinScore:      category.MinScore,
		Duration:      category.Duration,
		LimitQuestion: category.LimitQuestion,
		Message: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Category Updated!",
		},
	}
	dataCategory, _ := json.MarshalIndent(data, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataCategory)

}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
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

	cekIdCategory := controller.GetMaterialID(int(id), utils.Ctg)
	if cekIdCategory == 0 {
		res := response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Category Not Found!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	CekCategoryInQuestion := controller.CekMaterialInOtherRelation(int(id), "category_id", utils.Qst)
	if CekCategoryInQuestion == int(id) {
		res := response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Category Used In Question!",
		}
		data, _ := json.MarshalIndent(res, "", "\t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}
	sqlStatement := `DELETE FROM category WHERE id = $1`
	_, err = utils.DB.Exec(sqlStatement, id)
	row := utils.DB.QueryRow(sqlStatement, id)
	switch err := row.Scan(&category.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	res := response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Category Deleted!",
	}
	data, _ := json.MarshalIndent(res, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
