package main

import (
	"context"
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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderPage(w, "home", nil)
}

func (app *application) bookshelf(w http.ResponseWriter, r *http.Request) {
	// reading := []internals.Book{
	// 	{
	// 		Name:  "SQL Antipatterns, Volume 1",
	// 		Link:  "https://pragprog.com/titles/bksap1/sql-antipatterns-volume-1/",
	// 		Image: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1348552520l/7959038.jpg",
	// 	},
	// }

	// done := []internals.Book{
	// 	{
	// 		Name:  "The Art Of Postgres",
	// 		Link:  "https://theartofpostgresql.com/",
	// 		Image: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1567060454l/52755483._SX318_SY475_.jpg",
	// 	},
	// }

	// wantToRead := []internals.Book{
	// 	{
	// 		Name:  "Another Monster: investigate special",
	// 		Link:  "https://www.amazon.com/Another-Monster-investigative-special-%E3%82%82%E3%81%86%E3%81%B2%E3%81%A8%E3%81%A4%E3%81%AEMONSTER/dp/4091852793",
	// 		Image: "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1344467528l/6076366.jpg",
	// 	},
	// }
	ctx := context.Background()
	reading, err := app.queries.GetReadingBook(ctx)
	if err != nil {
		app.serverError(w, err)
	}
	done, err := app.queries.GetDoneBooks(ctx)
	if err != nil {
		app.serverError(w, err)
	}
	wantToRead, err := app.queries.GetWantToReadBooks(ctx)
	if err != nil {
		app.serverError(w, err)
	}

	books := internals.Books{
		Reading:  reading,
		Done:     done,
		WantRead: wantToRead,
	}

	renderPage(w, "bookshelf", books)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "about", nil)
}

func (app *application) quotes(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "quotes", nil)
}

func (app *application) anime(w http.ResponseWriter, r *http.Request) {
	url := "https://api.myanimelist.net/v2/users/UmaruKh/animelist?limit=500"
	var res internals.MyAnimeList

	cachedResult, found := cache.Get("anime")
	if found {
		res = cachedResult.(internals.MyAnimeList)
	} else {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			app.serverError(w, err)
			return
		}
		req.Header.Add("X-MAL-Client-ID", `6114d00ca681b7701d1e15fe11a4987e`)
		resp, err := client.Do(req)
		if err != nil {
			app.serverError(w, err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			app.serverError(w, err)
			return
		}
		cache.Add("anime", res, goCache.DefaultExpiration)
	}

	renderPage(w, "anime", res)
}

func (app *application) talks(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "talks", nil)
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "blog", nil)
}

func (app *application) openSource(w http.ResponseWriter, r *http.Request) {
	url := "https://codeberg.org/api/v1/users/omarkhatib/repos?limit=500"
	var res internals.GithubRepos

	cachedResult, found := cache.Get("open-source")
	if found {
		res = cachedResult.(internals.GithubRepos)
	} else {
		resp, err := http.Get(url)
		if err != nil {
			app.serverError(w, err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = json.Unmarshal(body, &res)
		if err != nil {
			app.serverError(w, err)
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
