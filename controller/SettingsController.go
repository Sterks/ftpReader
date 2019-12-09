package controller

import (
	"html/template"
	"log"
	"net/http"
)

// ViewData1 ...
type ViewData1 struct {
	Title string
}

// SettingsHandler ...
func SettingsHandler(w http.ResponseWriter, r *http.Request) {

	data := ViewData1{
		Title: "Settings",
	}

	tmpl, err := template.ParseFiles("view/Layout/Layout.html", "view/Body.html", "view/Profiles.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, data)
}
