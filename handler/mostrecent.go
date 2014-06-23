package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/dgonyeo/brandreth2.0/db"
)

func (h Handler) MostRecentTrip(w http.ResponseWriter, req *http.Request) {
	c := new(db.Controller)
	entries := c.GetLastTrip()
	log.Debug(strconv.Itoa(len(entries)))
	var models []*PersonEntry
	for _, entry := range entries {
		model := new(PersonEntry)
		model.Person = c.GetPerson(entry.UserId)
		model.Entry = entry
		models = append(models, model)
	}

	t, err := template.ParseFiles("templates/mostrecent.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, models)
	if err != nil {
		log.Fatal(err)
	}
}
