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
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Ищем по ID в бд (поиск по первичному  ключу так как используем gorm.Model)
	var task TaskServise.Task
	result := Database.DB.First(&task, id)
	if result.Error != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Декодируем тело запроса в мапу
	var update map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Обновляем Task декодированным запросом помещенным в мапу
	result = Database.DB.Model(&task).Updates(update)
	if result.Error != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// ID из URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Ищем задачу по ID
	var task TaskServise.Task
	result := Database.DB.First(&task, id)
	if result.Error != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Удаляем задачу
	result = Database.DB.Delete(&task)
	if result.Error != nil {
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
