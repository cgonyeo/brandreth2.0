package handler

import (
	"html/template"
	"net/http"

	"github.com/dgonyeo/brandreth2.0/config"
	"github.com/mholt/binding"
)

type SearchPage struct {
	PeopleEntries []*PersonEntry
	SearchQuery   string
}

func (sp SearchPage) IsActivePage(num int) bool {
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

	t, err := template.ParseFiles(
		config.Config.Templates.Path+"templates/search.tmpl",
		config.Config.Templates.Path+"templates/stuff.tmpl")
	if err != nil {
		log.Error("Error parsing the templates: %v", err)
		return
	}
	err = t.Execute(w, searchPage)
	if err != nil {
		log.Error("Error executing the templates: %v", err)
		return
	}
}
