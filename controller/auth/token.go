package auth

import (
	"errors"
	"fmt"
	"main/response"
	"main/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateToken(id int, role string) (response.Token, error) {
	if id == 0 {
		return response.Token{
			Token:  "undefined token",
			Name:   "user not valid!",
			UserID: 0,
		}, nil
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, _ := token.SignedString([]byte(os.Getenv("API_SECRET")))
	name := utils.GetAdminName(int(id), utils.Adm)
	return response.Token{
		Token:  tok,
		Name:   name,
		UserID: id,
	}, nil

}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(r *http.Request) (uint32, string, error) {

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, "", err
		}
		role := claims["role"]
		if str, ok := role.(string); ok {
			return uint32(id), str, nil
		} else {
			return 0, "", nil
		}
	}
	return 0, "", nil
}

func MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := TokenValid(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
