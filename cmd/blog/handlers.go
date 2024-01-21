package main

import (
	"log"
	"net/http"
	"text/template"
)

var _ Handlers = &handlers{}

type handlers struct {
}

type Handlers interface {
	Home(http.ResponseWriter, *http.Request)
	Contact(http.ResponseWriter, *http.Request)
	About(http.ResponseWriter, *http.Request)
}

func (h *handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	customTemplates := []string{
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
	w.Write([]byte("contact"))
}

func (h *handlers) About(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("about"))
}
