package controller

import (
	"database/sql"
	"fmt"
	"log"
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

func SignIn(email, password, sqlQuery string, id int) (response.Token, error) {

	var err error
	user := model.BaseUser{}
	log.Print("IDUSER:: ", id)
	log.Print("EMAIL:: ", email)
	log.Print("PASSWORD:: ", password)

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
	tok, _ := auth.CreateToken(id)
	return tok, nil
}
