package main

import (
	"fmt"
	"net/http"
	"os"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := fmt.Sprintf("%s - %s %s %s", ReadUserIP(r), r.Proto, r.Method, r.URL.RequestURI())

		f, err := os.OpenFile("ips.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			app.errorLog.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(log + "\n"); err != nil {
			app.errorLog.Println(err)
		}

		next.ServeHTTP(w, r)
	})
}
