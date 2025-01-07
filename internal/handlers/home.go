package handlers

import (
	"html/template"
	"net/http"
)

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates/home.html")
	if err != nil {
		http.Error(w, "Error loading homepage", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
