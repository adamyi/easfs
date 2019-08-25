package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2"
)

var (
	flagListenAddress string
	flagProd          bool
)

func GetEASFSPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Server", "easfs")

	ext := filepath.Ext(r.URL.Path)
	if ext == ".md" || ext == ".html" {
		http.Redirect(w, r, strings.TrimSuffix(r.URL.Path, ext), 301)
		return
	}

	// expiration := time.Now().Add(time.Hour)
	// cookie := http.Cookie{Name: "hl", Value: language, Expires: expiration}
	// http.SetCookie(w, &cookie)

	url := r.URL.Path
	if url == "/_s/getsuggestions" {
		if r.URL.Query().Get("c") == "2" {
			if r.URL.Query().Get("p") == "" {
				url = "/_suggestions"
			} else {
				url = filepath.Join("/", r.URL.Query().Get("p"), "/_suggestions")
			}
		} else {
			url = "/_empty_suggestions"
		}
	}

	var err error
	if IsDir("src/content/" + url) {
		// make sure that directory ends with a /
		if !strings.HasSuffix(url, "/") {
			http.Redirect(w, r, r.URL.Path+"/", 301)
			return
		}
		err = GetPage(w, url)
		if err.Error() != "file not found" {
			ReturnError(w, EASFSError{Code: http.StatusInternalServerError, Title: "500 Internal Server Error", Description: err.Error()})
			return
		}
		if err != nil {
			err = GetIndex(w, url)
		}
	} else {
		err = GetPage(w, url)
	}
	if err != nil {
		if err.Error() != "file not found" {
			ReturnError(w, EASFSError{Code: http.StatusInternalServerError, Title: "500 Internal Server Error", Description: err.Error()})
			return
		}
		red, err := GetRedirect(url)
		if err == nil {
			http.Redirect(w, r, red, 301)
		} else {
			ReturnError(w, EASFSError{Code: http.StatusNotFound, Title: "404 Not Found", Description: "The requested URL was not found on this server."})
		}
	}

	// fmt.Fprintf(w, "EASFS serving!\n")

}

func main() {
	flag.StringVar(&flagListenAddress, "listen_address", "0.0.0.0:80", "HTTP listen address")
	flag.BoolVar(&flagProd, "prod", false, "prod env")
	flag.Parse()
	pongo2.RegisterFilter("slugify", Slugify)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/_static/", http.StripPrefix("/_static/", fs))
	http.HandleFunc("/", GetEASFSPage)
	err := http.ListenAndServe(flagListenAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
