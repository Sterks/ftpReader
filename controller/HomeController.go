package controller

import (
	"html/template"
	"net/http"
)

// ViewData2 ...
type ViewData2 struct {
	Hello string
}

// HomeHandler ...
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	data := ViewData2{
		Hello: "Добрый День",
	}

	tmpl := template.Must(template.ParseFiles("view/Layout/Layout.html", "view/Home/Home.html"))

	tmpl.Execute(w, data)
}
