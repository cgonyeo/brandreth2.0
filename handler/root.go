package handler

import (
	"html/template"
	"net/http"
)

type RootPage struct {
    SearchQuery string
}

func (pp RootPage) IsActivePage(num int) bool {
    return false
}

func (h *Handler) Root(w http.ResponseWriter, req *http.Request) {
    model := new(RootPage)

	t, err := template.ParseFiles("templates/root.tmpl", "templates/stuff.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, model)
	if err != nil {
		log.Fatal(err)
	}
}
