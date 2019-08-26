package main

import (
	"net/http"
	"text/template"
)

type EASFSError struct {
	Title       string
	Description string
	Code        int
}

func ReturnError(w http.ResponseWriter, err EASFSError) {
	w.WriteHeader(err.Code)
	tmpl := template.Must(template.ParseFiles("templates/error.html"))
	tmpl.Execute(w, err)
}
