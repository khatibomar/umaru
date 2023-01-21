package main

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type ViewData struct {
	Year string
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderPage(w, "home")
}

func readings(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "readings")
}

func about(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "about")
}

func quotes(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "quotes")
}

func snippets(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "snippets")
}

func talks(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "talks")
}

func renderPage(w http.ResponseWriter, page string) {
	files := []string{
		"./templates/base.html",
		"./templates/" + page + ".html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	year := ViewData{
		Year: strconv.Itoa(time.Now().Year()),
	}
	err = ts.ExecuteTemplate(w, "base", year)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/readings", readings)
	mux.HandleFunc("/talks", talks)
	mux.HandleFunc("/quotes", quotes)
	mux.HandleFunc("/snippets", snippets)
	mux.HandleFunc("/about", about)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
