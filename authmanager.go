package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		pw := r.FormValue("password")

		items := struct {
			Username string
			Password string
		}{
			Username: un,
			Password: pw,
		}

		am.loginTemplate.Execute(w, items)
		return
	}

	items := struct {
		Country string
		City    string
	}{
		Country: "Australia",
		City:    "Paris",
	}

	am.loginTemplate.Execute(w, items)
}

func adminHandlerTwo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It's adminHandlerTwo , Hello, %q", r.URL.Path[1:])
}
