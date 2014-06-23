package handler

import (
	"html/template"
	"net/http"

	"github.com/dgonyeo/brandreth2.0/db"
	"github.com/mholt/binding"
)

type SearchParams struct {
	Search string
}

func (sp *SearchParams) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&sp.Search: "search",
	}
}

func (h Handler) Search(w http.ResponseWriter, req *http.Request) {
	searchParams := new(SearchParams)
	errs := binding.Bind(req, searchParams)
	if errs.Handle(w) {
		log.Error("Error with binding")
		return
	}

	c := new(db.Controller)
	entries := c.SearchForTrips(searchParams.Search)
    var peopleEntries []*PersonEntry
	for _, entry := range entries {
		pe := new(PersonEntry)
		pe.Person = c.GetPerson(entry.UserId)
		pe.Entry = entry
		peopleEntries = append(peopleEntries, pe)
	}

	t, err := template.ParseFiles("templates/search.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, peopleEntries)
	if err != nil {
		log.Fatal(err)
	}
}
