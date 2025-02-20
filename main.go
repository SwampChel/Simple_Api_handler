package main

import (
	"ApiHandler/Database"
	"ApiHandler/TaskServise"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var tasks []TaskServise.Task

	if err := Database.DB.Find(&tasks).Error; err != nil { // извлекаем
		http.Error(w, "Fail", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(tasks) //кодируем в формат json
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var task TaskServise.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "It's Over", http.StatusBadRequest)
		return
	}
	if err := Database.DB.Create(&task).Error; err != nil { // cоздаем запись в базе
		http.Error(w, "It's Over ", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(task) // кодируем в формат json ответ
}

func PatchHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// ID из URL
	idVars := mux.Vars(r)
	id, err := strconv.Atoi(idVars["id"])
	if err != nil {
		http.Error(w, "you are invalid ", http.StatusBadRequest)
		return
	}

	var task TaskServise.Task

	// Декодируем тело запроса в мапу

	var update map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Обновляем Task декодированным запросом помещенным в мапу
	newTask := Database.DB.Model(&task).Where("id = ?", id).Updates(update)
	if newTask.Error != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// ID из URL
	idVars := mux.Vars(r)
	id, err := strconv.Atoi(idVars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Удаляем задачу

	var task TaskServise.Task

	delTask := Database.DB.Delete(&task, id)
	if delTask.Error != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := mux.NewRouter()

	Database.InitDB()

	Database.DB.AutoMigrate(&TaskServise.Task{})

	router.HandleFunc("/api/get", GetHandler).Methods("GET")
	router.HandleFunc("/api/post", CreateHandler).Methods("POST")
	router.HandleFunc("/api/patch/{id}", PatchHandler).Methods("PATCH")
	router.HandleFunc("/api/delete/{id}", DeleteHandler).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
