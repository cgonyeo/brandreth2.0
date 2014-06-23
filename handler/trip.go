package handler

import (
	"html/template"
	"net/http"

	"github.com/dgonyeo/brandreth2.0/db"
	"github.com/mholt/binding"
)

type TripPage struct {
	TripInfo      *db.Entry
	PeopleEntries []*PersonEntry
}

type TripParams struct {
	TripId string
}

func (pp *TripParams) FieldMap() binding.FieldMap {
	return binding.FieldMap{
        &pp.TripId: binding.Field{
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

	c := new(db.Controller)
	entries := c.GetTripsEntries(tripParams.TripId)
	model := new(TripPage)
	for _, entry := range entries {
		if model.TripInfo == nil {
			model.TripInfo = entry
		}
		pe := new(PersonEntry)
		pe.Person = c.GetPerson(entry.UserId)
		pe.Entry = entry
		model.PeopleEntries = append(model.PeopleEntries, pe)
	}

	t, err := template.ParseFiles("templates/trip.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, model)
	if err != nil {
		log.Fatal(err)
	}
}
