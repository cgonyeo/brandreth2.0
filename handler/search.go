package handler

import (
	"html/template"
	"net/http"

	"github.com/mholt/binding"
)

type SearchPage struct {
	PeopleEntries []*PersonEntry
	SearchQuery   string
}

func (sp SearchPage) IsActivePage(num int) bool {
	log.Debug("SearchPage")
	return false
}

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

	entries := h.c.SearchForTrips(searchParams.Search)
	searchPage := new(SearchPage)
	for _, entry := range entries {
		pe := new(PersonEntry)
		pe.Person = h.c.GetPerson(entry.UserId)
		pe.Entry = entry
		searchPage.PeopleEntries = append(searchPage.PeopleEntries, pe)
	}
	searchPage.SearchQuery = searchParams.Search

	t, err := template.ParseFiles("templates/search.tmpl", "templates/stuff.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, searchPage)
	if err != nil {
		log.Fatal(err)
	}
}
