package main

import (
	"ApiHandler/Database"
	"ApiHandler/TaskServise"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {

	var tasks []TaskServise.Task

	if err := Database.DB.Find(&tasks).Error; err != nil { // cоздаем запись в базе
		http.Error(w, "Fail", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json") //явно указываем формат json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks) //кодируем в формат json
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {

	var task TaskServise.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "It's Over", http.StatusBadRequest)
		return
	}
	if err := Database.DB.Create(&task).Error; err != nil { //Извлекаем
		http.Error(w, "It's UberOver ", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json") //явно указываем формат json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task) // кодируем в формат json
}

func main() {
	router := mux.NewRouter()

	Database.InitDB()

	Database.DB.AutoMigrate(&TaskServise.Task{})

	router.HandleFunc("/api/task", GetMessages).Methods("GET")
	router.HandleFunc("/api/task", CreateMessage).Methods("POST")

	http.ListenAndServe(":8080", router)
}
