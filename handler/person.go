package handler

import (
	"html/template"
	"net/http"

	"github.com/mholt/binding"

	"github.com/dgonyeo/brandreth2.0/db"
)

type PersonPage struct {
	Person      *db.Person
	Entries     []*db.Entry
	Months      []string
	Trips       []int
	SearchQuery string
}

func (pp PersonPage) IsActivePage(num int) bool {
	return false
}

type PersonParams struct {
	UserId string
}

func (pp *PersonParams) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&pp.UserId: "user_id",
	}
}

func (h *Handler) Person(w http.ResponseWriter, req *http.Request) {
	personParams := new(PersonParams)
	errs := binding.Bind(req, personParams)
	if errs.Handle(w) {
		log.Error("Error with binding")
		return
	}

	person := h.c.GetPerson(personParams.UserId)
	entries := h.c.GetPersonsEntries(personParams.UserId)
	model := new(PersonPage)
	model.Person = person
	model.Entries = entries
	model.Months, model.Trips = h.c.GetMonthCountForPerson(person.UserId)

	t, err := template.ParseFiles("templates/person.tmpl", "templates/stuff.tmpl")
	if err != nil {
		return
	}
	err = t.Execute(w, model)
	if err != nil {
		return
	}
}
