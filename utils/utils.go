package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"main/configuration"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	DB = configuration.OpenConnection()
)

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("failed to hash password")
	}
	return string(hash)
}

func VerifyPassword(hashedPassword, password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return
}

func FormatError(err string) error {

	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}

	if strings.Contains(err, "nik") {
		return errors.New("NIK Already Taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}

	return errors.New("Incorrect Details")
}

func GetName(id int, tbl string) string {

	var name string
	sqlQuery := `SELECT name FROM ` + tbl + ` WHERE id = $1`

	row := DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(&name); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	}

	return name
}
