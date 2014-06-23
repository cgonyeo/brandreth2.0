package handler

import (
	"html/template"
	"net/http"

	"github.com/dgonyeo/brandreth2.0/db"
	"github.com/mholt/binding"
)

type PersonPage struct {
	Person  *db.Person
	Entries []*db.Entry
}

type PersonParams struct {
	UserId string
}

func (pp *PersonParams) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&pp.UserId: "user_id",
	}
}

func (h Handler) Person(w http.ResponseWriter, req *http.Request) {
	personParams := new(PersonParams)
	errs := binding.Bind(req, personParams)
	if errs.Handle(w) {
		log.Error("Error with binding")
		return
	}

	c := new(db.Controller)
	person := c.GetPerson(personParams.UserId)
	entries := c.GetPersonsEntries(personParams.UserId)
	model := new(PersonPage)
	model.Person = person
	model.Entries = entries

	t, err := template.ParseFiles("templates/person.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, model)
	if err != nil {
		log.Fatal(err)
	}
}
