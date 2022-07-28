package controller

import (
	"database/sql"
	"fmt"
	"main/controller/auth"
	"main/model"
	"main/query"
	"main/response"
	"main/utils"

	"golang.org/x/crypto/bcrypt"
)

func CekUser(email, sqlQuery string) response.UserResponse {

	var user response.UserResponse

	row := utils.DB.QueryRow(sqlQuery, email)
	switch err := row.Scan(&user.Name, &user.Email); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return user
}

func GetUserID(email, sqlQuery string) int {

	var user model.BaseUser

	row := utils.DB.QueryRow(sqlQuery, email)
	switch err := row.Scan(&user.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return user.ID
}

func SignIn(email, password, sqlQuery, role string, id int) (response.Token, error) {

	var err error
	user := model.BaseUser{}

	row := utils.DB.QueryRow(sqlQuery, email)
	err = row.Err()
	if err != nil {
		return response.Token{
			Token:  "",
			Name:   "",
			UserID: 0,
		}, err
	}
	err = utils.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return response.Token{
			Token:  "",
			Name:   "",
			UserID: 0,
		}, err
	}
	tok, _ := auth.CreateToken(id, role)
	return tok, nil
}

func GetMaterialID(id int, tbl string) int {

	var idMaterial response.IdResponse
	sqlQuery := query.SqlGetMaterialID(tbl)

	row := utils.DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(&idMaterial.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return idMaterial.ID
}

func CekMaterialInOtherRelation(id int, column, tbl string) int {

	var idMaterial response.IdResponse
	sqlQuery := query.SqlCekMaterialInOtherRelation(id, column, tbl)

	row := utils.DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(&idMaterial.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return idMaterial.ID
}

func GetCategoryIDFromQuestion(id int) int {

	var idMaterial response.IdResponse
	sqlQuery := query.SqlGetCategoryIDFromQuestion()

	row := utils.DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(&idMaterial.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return idMaterial.ID
}

func GetQuestionIDFromAnswer(id int) int {

	var idMaterial response.IdResponse
	sqlQuery := query.SqlGetQuestionIDFromAnswer()

	row := utils.DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(&idMaterial.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return idMaterial.ID
}
