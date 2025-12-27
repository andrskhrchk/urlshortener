package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var storage *Storage

func StartServer(s *Storage) {
	storage = s
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("/shorten", handleShorten)

	mux.HandleFunc("/r/", handleRedirect)

	fmt.Println("Server listening to :8080")
	http.ListenAndServe(":8080", mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)

}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	longURL := r.FormValue("long_url")
	fmt.Printf("Ссылка: %s\n", longURL)

	id, err := storage.SaveURL(longURL)
	if err != nil {
		log.Printf("Ошибка записи в БД: %v", err)
		http.Error(w, "Ошибка сохранения", 500)
		return
	}

	code := EncodeBase62(id)

	if storage.AddShortCode(id, code) != nil {
		log.Printf("Ошибка при обновлении short_code в БД: %v", err)
		http.Error(w, "Не удалось завершить создание ссылки", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Сокращённая ссылка: http://localhost:8080/r/%s\n", code)

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]

	longURL, err := storage.GetLongURL(code)
	if err != nil {
		http.Error(w, "Ссылка не найдена", 500)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
