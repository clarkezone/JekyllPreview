package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AuthManager is a web server for managing signin
type AuthManager struct {
	authmux       *http.ServeMux
	loginTemplate *template.Template
}

// NewAuthManager creates a new instance
func NewAuthManager() *AuthManager {
	httpMan := &AuthManager{}
	httpMan.parseTemplates()
	adminMux := http.NewServeMux()

	adminMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	adminMux.HandleFunc("/login", httpMan.adminHandlerOne)
	adminMux.HandleFunc("/index", httpMan.adminHandlerTwo)
	httpMan.authmux = adminMux
	return httpMan
}

func (am *AuthManager) parseTemplates() {
	t, err := template.ParseFiles("./htmltemplates/login.html")

	if err != nil {
		log.Fatalf("Failed to load template: %v\n", err.Error())
	}
	am.loginTemplate = t
}

func (am *AuthManager) adminHandlerOne(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		re := r.URL.Query().Get("r")

		items := struct {
			Username    string
			Password    string
			RedirectURL string
			LoginFailed bool
		}{
			Username:    "",
			Password:    "",
			RedirectURL: re,
			LoginFailed: false,
		}

		am.loginTemplate.Execute(w, items)
		return
	}

	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		pw := r.FormValue("password")
		redir := r.FormValue("redirectto")

		re, err := url.QueryUnescape(redir)

		if am.validateUser(un, pw) {
			am.setCookie(w, r.Host)

			if err == nil && re != "" {
				http.Redirect(w, r, re, http.StatusSeeOther)
				return
			}
		} else {
			items := struct {
				Username    string
				Password    string
				RedirectURL string
				LoginFailed bool
			}{
				Username:    un,
				Password:    "",
				RedirectURL: re,
				LoginFailed: true,
			}
			am.loginTemplate.Execute(w, items)
			return
		}
	}

}

func (am *AuthManager) adminHandlerTwo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "handler two")
}

//TODO wrap this in an interface
func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkCookie(r) {
			host := authman.stripSubdomain(r)

			var proto string
			if proto = "http://"; r.TLS != nil {
				proto = "https://"
			}

			newURI := proto + host + "/login?r=" + url.QueryEscape(proto+r.Host+r.RequestURI)
			http.Redirect(w, r, newURI, http.StatusSeeOther)
			return
		}
		fn(w, r)
	}
}

func (am *AuthManager) stripSubdomain(r *http.Request) string {
	domainParts := strings.Split(r.Host, ".")
	if len(domainParts) >= 3 && domainParts[0] != "preview" {
		newhost := strings.Replace(r.Host, domainParts[0], "", 1)
		newhost = strings.TrimPrefix(newhost, ".")
		return newhost
	}
	return r.Host + r.RequestURI
}

func (am *AuthManager) validateUser(u, p string) bool {
	if u == "clearsky" && p == "froyo2020" {
		return true
	}
	return false
}

func (am *AuthManager) setCookie(rw http.ResponseWriter, domain string) {
	//TODO hard coded
	expire := time.Now().Add(5 * time.Minute)
	// in order for subdomains to see the cookies correctly I found that only "preview.localhost" works, nothing else does.  As a consequence
	// if you don't build a domain manually which strips off the port, no dice.
	d2 := strings.Replace(domain, ":8085", "", 1)
	cookie := http.Cookie{Name: "session", Value: "loggied in", Expires: expire, Domain: d2}
	http.SetCookie(rw, &cookie)
}

func checkCookie(r *http.Request) bool {
	session, err := r.Cookie("session")
	if session == nil || err != nil {
		return false
	}
	return true
}
