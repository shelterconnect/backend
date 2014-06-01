package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/zachlatta/shelterconnect/database"
	"github.com/zachlatta/shelterconnect/handler"
)

func httpLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
		log.Printf("Completed in %s", time.Now().Sub(start).String())
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err := database.Init("postgres",
		os.ExpandEnv("postgres://docker:docker@$DB_1_PORT_5432_TCP_ADDR/docker"))
	if err != nil {
		panic(err)
	}
	defer database.Close()

	r := mux.NewRouter()

	r.Handle("/organizations",
		handler.AppHandler(handler.CreateOrganization)).Methods("POST")
	r.Handle("/organizations",
		handler.AppHandler(handler.GetAllOrganizations)).Methods("GET")
	r.Handle("/organizations/{id}",
		handler.AppHandler(handler.GetOrganization)).Methods("GET")
	r.Handle("/organizations/authenticate",
		handler.AppHandler(handler.AuthenticateOrganization)).Methods("POST")

	r.Handle("/shelters",
		handler.AppHandler(handler.GetAllShelters)).Methods("GET")
	r.Handle("/shelters/{id}",
		handler.AppHandler(handler.GetShelter)).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":"+port, httpLog(http.DefaultServeMux)))
}
