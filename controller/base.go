package controller

import (
	"database/sql"
	"fmt"
	"main/configuration"
	"main/response"
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
