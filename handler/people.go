package handler

import (
	"html/template"
	"net/http"

	"github.com/dgonyeo/brandreth2.0/db"
)

type PeoplePage struct {
    People []*db.Person
    SearchQuery string
}

func (pp PeoplePage) IsActivePage(num int) bool {
    return num == 2
}

func (h *Handler) People(w http.ResponseWriter, req *http.Request) {
	c := new(db.Controller)
	people := c.GetPeople()

    t, err := template.ParseFiles("templates/people.tmpl", "templates/stuff.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, PeoplePage{people, ""})
	if err != nil {
        log.Fatal("People: %v", err)
	}
}
