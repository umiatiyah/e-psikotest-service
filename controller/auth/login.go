package auth

import (
	"main/configuration"
	"main/model"
	"main/response"
	"main/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	db = configuration.OpenConnection()
)

func SignIn(email, password, sqlQuery string) (response.Token, error) {

	var err error
	user := model.UserAdmin{}

	row := db.QueryRow(sqlQuery, email)
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
	tok, _ := CreateToken(user.ID)
	return tok, nil
}
