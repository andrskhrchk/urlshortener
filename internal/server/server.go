package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"example.com/internal/storage"
)

type Server struct {
	repo storage.URLRepository
}

func NewServer(repo storage.URLRepository) *Server {
	return &Server{repo: repo}
}

func (srv *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", srv.handleRoot)

	mux.HandleFunc("/shorten", srv.handleShorten)

	mux.HandleFunc("/r/", srv.handleRedirect)

	fmt.Println("Server listening to :8080")
	http.ListenAndServe(":8080", mux)
}

func (srv *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)

}

func (srv *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	longURL := r.FormValue("long_url")
	fmt.Printf("Ссылка: %s\n", longURL)

	code, err := srv.repo.SaveURLShort(longURL)
	if err != nil {
		log.Printf("Ошибка записи в БД: %v", err)
		http.Error(w, "Ошибка сохранения", 500)
		return
	}

	fmt.Fprintf(w, "Сокращённая ссылка: http://localhost:8080/r/%s\n", code)

}

func (srv *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]

	longURL, err := srv.repo.GetLongURL(code)
	if err != nil {
		http.Error(w, "Ссылка не найдена", 500)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
