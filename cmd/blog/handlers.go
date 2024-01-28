package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"log"
	"net/http"
	"text/template"
)

var _ Handlers = &handlers{}

type handlers struct{}

type Handlers interface {
	Home(http.ResponseWriter, *http.Request)
	Contact(http.ResponseWriter, *http.Request)
	Blog(http.ResponseWriter, *http.Request)
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// Подключение к БД
func connectToDB(cfg DBConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}

func (h *handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	customTemplates := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
	}

	ts, err := template.ParseFiles(customTemplates...)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}

}

func (h *handlers) Contact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "contact")
}

type Post struct {
	Author      string
	H1          string
	Description string
	Text        string
	DateAdded   string
	ImgUrl      string
	Views       int
	Comments    int
	Likes       int
}

func (h *handlers) Blog(w http.ResponseWriter, r *http.Request) {
	cfg := DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "seoSuperUser",
		Password: "ronin95PG",
		DBName:   "blog",
	}

	db, err := connectToDB(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT author,h1,description,text,date_added,img_url,views,comments,likes FROM post")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	customTemplates := []string{
		"./ui/html/blog.page.html",
		"./ui/html/base.layout.html",
	}

	ts, err := template.ParseFiles(customTemplates...)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var post []Post

	for rows.Next() {
		var p Post
		if err = rows.Scan(&p.Author, &p.H1, &p.Description, &p.Text, &p.DateAdded, &p.ImgUrl, &p.Views, &p.Comments, &p.Likes); err != nil {
			log.Println(err)
			return
		}
		post = append(post, p)
	}

	// Проверяем наличие ошибок при переборе строк.
	if err = rows.Err(); err != nil {
		log.Println(err)
	}

	if err = ts.Execute(w, post); err != nil {
		log.Println(err)
	}

}
