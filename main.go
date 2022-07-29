package main

import (
	"log"
	"main/server"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	r := mux.NewRouter()

	server.InitializeRoute(r)
	handle := cors.AllowAll().Handler(r)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handle))

}
