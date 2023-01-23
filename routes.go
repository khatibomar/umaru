package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/open-source", app.openSource)
	mux.HandleFunc("/bookshelf", app.bookshelf)
	mux.HandleFunc("/talks", app.talks)
	mux.HandleFunc("/quotes", app.quotes)
	mux.HandleFunc("/anime", app.anime)
	mux.HandleFunc("/blog", app.blog)
	mux.HandleFunc("/about", app.about)

	return mux
}
