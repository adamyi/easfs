package main

import (
	"context"
	"flag"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/flosch/pongo2"
	"github.com/golang/glog"
)

var (
	flagListenAddress    string
	flagSSLListenAddress string
	flagSSLCert          string
	flagSSLKey           string
	flagDomain           string
	flagProd             bool
	flagSitePath         string
)

func EASFSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Server", "easfs")
	if r.Host != flagDomain {
		RedirectDomain(w, r)
		return
	}

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
	} else if strings.Contains(url, "/_") {
		ReturnError(w, EASFSError{Code: http.StatusNotFound, Title: "404 Not Found", Description: "The requested URL was not found on this server."})
		return
	}

	var err error
	if IsDir(flagSitePath + url) {
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

func RedirectSSL(rsp http.ResponseWriter, req *http.Request) {
	rsp.Header().Set("X-Frame-Options", "SAMEORIGIN")
	rsp.Header().Set("Server", "easfs")
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(rsp, req, target,
		http.StatusPermanentRedirect)
}

func RedirectDomain(rsp http.ResponseWriter, req *http.Request) {
	target := "https://" + flagDomain + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(rsp, req, target,
		http.StatusTemporaryRedirect)
}

var siteStatusHTML = `Site at {{.SitePath}}<br>Site version {{.SiteVersion}}`
var SiteVersion = "Unknown"

func siteStatus(context.Context) interface{} {
	data := struct {
		SitePath    string
		SiteVersion string
	}{
		SitePath:    flagSitePath,
		SiteVersion: SiteVersion,
	}
	return data
}

func main() {
	AddStatusPart("Site Status", siteStatusHTML, siteStatus)
	flag.StringVar(&flagListenAddress, "listen", "0.0.0.0:80", "HTTP listen address")
	flag.StringVar(&flagSSLListenAddress, "slisten", "0.0.0.0:443", "HTTPS listen address")
	flag.StringVar(&flagSSLCert, "cert", "cert.pem", "HTTPS cert")
	flag.StringVar(&flagSSLKey, "key", "key.pem", "HTTPS key")
	flag.StringVar(&flagDomain, "domain", "", "Site domain")
	flag.StringVar(&flagSitePath, "site", "site/content/", "Path to site content")
	flag.BoolVar(&flagProd, "prod", false, "prod env")
	flag.Parse()
	cmd := exec.Command("git", "describe", "--always", "--match", "v[0-9].*", "--dirty")
	cmd.Dir = flagSitePath
	out, err := cmd.Output()
	if err == nil {
		SiteVersion = strings.TrimSpace(string(out))
	}
	pongo2.RegisterFilter("slugify", Slugify)
	go func() {
		httperr := http.ListenAndServe(flagListenAddress, http.HandlerFunc(RedirectSSL))
		if httperr != nil {
			glog.Fatal("ListenAndServe: ", httperr)
		}
	}()
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/_static/", gziphandler.GzipHandler(http.StripPrefix("/_static/", fs)))
	mux.Handle("/", gziphandler.GzipHandler(http.HandlerFunc(EASFSHandler)))
	err = http.ListenAndServeTLS(flagSSLListenAddress, flagSSLCert, flagSSLKey, mux)
	if err != nil {
		glog.Fatal("ListenAndServeTLS: ", err)
	}
}
