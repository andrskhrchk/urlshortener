package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"example.com/internal/storage"
)

type Server struct {
	repo storage.URLRepository
}

func (srv *Server) renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmplPath := filepath.Join("web", "templates", tmplName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Ошибка: файл %s не найден", tmplName)
		http.Error(w, "Ошибка страницы", 500)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Ошибка выполнения шаблона: %s: %v", tmplName, err)
		http.Error(w, "Ошибка при генерации страницы", 500)
	}
}

func NewServer(repo storage.URLRepository) *Server {
	return &Server{repo: repo}
}

func (srv *Server) Start() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", srv.handleRoot)

	mux.HandleFunc("/shorten", srv.handleShorten)

	mux.HandleFunc("/r/", srv.handleRedirect)

	mux.HandleFunc("/login", srv.handleLogin)

	mux.HandleFunc("/register", srv.handleRegister)

	fmt.Println("Server listening to :8080")
	http.ListenAndServe(":8080", mux)
}

func (srv *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	srv.renderTemplate(w, "index.html", nil)
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

	shortUrl := fmt.Sprintf("http://localhost:8080/r/%s", code)

	data := map[string]interface{}{
		"ShortURL": shortUrl,
	}

	srv.renderTemplate(w, "index.html", data)

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

func (srv *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	srv.renderTemplate(w, "login.html", nil)
}

func (srv *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	srv.renderTemplate(w, "register.html", nil)
}
