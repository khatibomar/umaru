package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"codeberg.org/omarkhatib/umaru/internals"
	goCache "github.com/patrickmn/go-cache"
)

var cache = goCache.New(12*time.Hour, 24*time.Hour)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderPage(w, "home", nil)
}

func bookshelf(w http.ResponseWriter, r *http.Request) {
	reading := []internals.Book{
		{
			Name:  "SQL Antipatterns, Volume 1",
			Link:  "https://pragprog.com/titles/bksap1/sql-antipatterns-volume-1/",
			Image: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1348552520l/7959038.jpg",
		},
	}

	done := []internals.Book{
		{
			Name:  "The Art Of Postgres",
			Link:  "https://theartofpostgresql.com/",
			Image: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1567060454l/52755483._SX318_SY475_.jpg",
		},
	}

	wantToRead := []internals.Book{
		{
			Name:  "Another Monster: investigate special",
			Link:  "https://www.amazon.com/Another-Monster-investigative-special-%E3%82%82%E3%81%86%E3%81%B2%E3%81%A8%E3%81%A4%E3%81%AEMONSTER/dp/4091852793",
			Image: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1344467528l/6076366.jpg",
		},
	}

	books := internals.Books{
		Reading:  reading,
		Done:     done,
		WantRead: wantToRead,
	}

	renderPage(w, "bookshelf", books)
}

func about(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "about", nil)
}

func quotes(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "quotes", nil)
}

func anime(w http.ResponseWriter, r *http.Request) {
	url := "https://api.myanimelist.net/v2/users/UmaruKh/animelist?limit=500"
	var res internals.MyAnimeList

	cachedResult, found := cache.Get("anime")
	if found {
		res = cachedResult.(internals.MyAnimeList)
	} else {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		req.Header.Add("X-MAL-Client-ID", `6114d00ca681b7701d1e15fe11a4987e`)
		resp, err := client.Do(req)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		cache.Add("anime", res, goCache.DefaultExpiration)
	}

	renderPage(w, "anime", res)
}

func talks(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "talks", nil)
}

func blog(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "blog", nil)
}

func openSource(w http.ResponseWriter, r *http.Request) {
	url := "https://codeberg.org/api/v1/users/omarkhatib/repos?limit=500"
	var res internals.GithubRepos

	cachedResult, found := cache.Get("open-source")
	if found {
		res = cachedResult.(internals.GithubRepos)
	} else {
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		cache.Add("open-source", res, goCache.DefaultExpiration)
	}

	renderPage(w, "opensource", res)
}

func renderPage(w http.ResponseWriter, page string, data any) {
	files := []string{
		"./templates/base.html",
		"./templates/partials/nav.html",
		"./templates/pages/" + page + ".html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	bag := internals.ViewData{
		Data: data,
		Year: strconv.Itoa(time.Now().Year()),
	}
	err = ts.ExecuteTemplate(w, "base", bag)
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
	mux.HandleFunc("/open-source", openSource)
	mux.HandleFunc("/bookshelf", bookshelf)
	mux.HandleFunc("/talks", talks)
	mux.HandleFunc("/quotes", quotes)
	mux.HandleFunc("/anime", anime)
	mux.HandleFunc("/blog", blog)
	mux.HandleFunc("/about", about)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
