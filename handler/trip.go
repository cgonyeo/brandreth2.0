package handler

import (
	"html/template"
	"net/http"

	"github.com/mholt/binding"

	"github.com/dgonyeo/brandreth2.0/config"
	"github.com/dgonyeo/brandreth2.0/db"
)

type TripPage struct {
	TripInfo      *db.Entry
	PeopleEntries []*PersonEntry
	SearchQuery   string
}

func (tp TripPage) IsActivePage(num int) bool {
	return false
}

type TripParams struct {
	TripId string
}

func (tp *TripParams) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&tp.TripId: binding.Field{
			Form:     "trip_id",
			Required: true,
		},
	}
}

func (h Handler) Trip(w http.ResponseWriter, req *http.Request) {
	tripParams := new(TripParams)
	errs := binding.Bind(req, tripParams)
	if errs.Handle(w) {
		log.Error("Error with binding")
		return
	}

	entries := h.c.GetTripsEntries(tripParams.TripId)
	model := new(TripPage)
	for _, entry := range entries {
		if model.TripInfo == nil {
			model.TripInfo = entry
		}
		pe := new(PersonEntry)
		pe.Person = h.c.GetPerson(entry.UserId)
		pe.Entry = entry
		model.PeopleEntries = append(model.PeopleEntries, pe)
	}

	t, err := template.ParseFiles(
		config.Config.Templates.Path+"templates/trip.tmpl",
		config.Config.Templates.Path+"templates/stuff.tmpl")
	if err != nil {
		log.Error("Error parsing the templates: %v", err)
		return
	}
	err = t.Execute(w, model)
	if err != nil {
		log.Error("Error executing the templates: %v", err)
		return
	}
}
