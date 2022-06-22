package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type IdResponse struct {
	ID int `json:"id"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}

type Token struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}
