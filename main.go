package main

import (
	"log"
	"main/server"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	server.InitializeRoute(r)

	log.Fatal(http.ListenAndServe(":8080", r))

}
