package controller

import (
	"database/sql"
	"fmt"
	"main/configuration"
	"main/controller/auth"
	"main/model"
	"main/response"
	"main/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	DB = configuration.OpenConnection()
)

func CekUser(email, sqlQuery string) response.UserResponse {

	var user response.UserResponse

	row := DB.QueryRow(sqlQuery, email)
	switch err := row.Scan(&user.Name, &user.Email); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return user
}

func SqlQueryCek(tbl string) string {
	return `SELECT name, email FROM ` + tbl + ` WHERE email = $1`
}

func GetUserID(email, sqlQuery string) int {

	var user model.BaseUser

	row := DB.QueryRow(sqlQuery, email)
	switch err := row.Scan(&user.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return user.ID
}

func SqlGetID(tbl string) string {
	return `SELECT id FROM ` + tbl + ` WHERE email = $1`
}

func SignIn(email, password, sqlQuery, role string, id int) (response.Token, error) {

	var err error
	user := model.BaseUser{}

	row := DB.QueryRow(sqlQuery, email)
	err = row.Err()
	if err != nil {
		return response.Token{
			Token: "",
		}, err
	}
	err = utils.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return response.Token{
			Token: "",
		}, err
	}
	tok, _ := auth.CreateToken(id, role)
	return tok, nil
}

func GetAdminName(id int) string {

	var admin model.Admin
	sqlQuery := `SELECT name FROM admin WHERE id = $1`

	row := DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(&admin.Name); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return admin.Name
}

func GetMaterialID(id int, tbl string) int {

	var idMaterial response.IdResponse
	sqlQuery := `SELECT id FROM ` + tbl + ` WHERE id = $1`

	row := DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(&idMaterial.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return idMaterial.ID
}

func CekMaterialInOtherRelation(id int, column, tbl string) int {

	var idMaterial response.IdResponse
	sqlQuery := `SELECT ` + column + ` FROM ` + tbl + ` WHERE id in ($1)`

	row := DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(&idMaterial.ID); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return idMaterial.ID
}
