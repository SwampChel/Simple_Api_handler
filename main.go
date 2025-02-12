package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var task string

type requestBody struct {
	Message string `json:"message"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "oh hi, %s", task)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var body requestBody

	// Декодируем json из тела запроса в структуру
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "It's Over", http.StatusBadRequest)
		return
	}
	// если ключ json соответствует то все хорошо приравниваем запрос к task
	task = body.Message

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Are you winning son ?")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/hello", GetHandler).Methods("GET")
	router.HandleFunc("/api/post", PostHandler).Methods("POST")

	http.ListenAndServe(":8080", router)
}
