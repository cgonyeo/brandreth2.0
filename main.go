package main

import (
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	golog "github.com/op/go-logging"

	"github.com/dgonyeo/brandreth2.0/db"
	"github.com/dgonyeo/brandreth2.0/handler"
	"github.com/dgonyeo/brandreth2.0/importer"
)

var log = golog.MustGetLogger("main")

func main() {
	webapp()
}

func importData() {
	controller := new(db.Controller)
	controller.CreateTables()

	importer.Run("guestbook.csv", controller)
}

func webapp() {
	h := new(handler.Handler)
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.MostRecentTrip)
	mux.HandleFunc("/person", h.Person)
	mux.HandleFunc("/trip", h.Trip)
	mux.HandleFunc("/search", h.Search)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3001")
}

func test() {
	controller := new(db.Controller)
	log.Debug("session acquired")

	controller.CreateTables()
	log.Debug("tables created")

	person := db.Person{"", "Derek Gonyeo", "dgonz", "CSH"}
	user_id := controller.AddPerson(&person)
	log.Debug("person created")

	entry1 := new(db.Entry)
	entry1.TripId = db.GetUniqueId()
	entry1.UserId = user_id
	entry1.TripReason = "To see some awesome trees"
	entry1.DateStart = time.Date(2014, 6, 20, 0, 0, 0, 0, new(time.Location))
	entry1.DateEnd = time.Date(2014, 6, 22, 0, 0, 0, 0, new(time.Location))
	entry1.Entry = "I saw normal trees /and/ Christmas trees"
	entry1.Book = 3
	controller.AddEntry(entry1)
	log.Debug("entry1 added")

	log.Debug(controller.GetPerson(user_id).String())
	log.Debug("Person retrieved")

	log.Debug(controller.GetEntry(user_id, entry1.TripId).String())
	log.Debug("Entry retrieved")

	entry2 := new(db.Entry)
	entry2.TripId = db.GetUniqueId()
	entry2.UserId = user_id
	entry2.TripReason = "To get in a canoe"
	entry2.DateStart = time.Date(2014, 6, 27, 0, 0, 0, 0, new(time.Location))
	entry2.DateEnd = time.Date(2014, 6, 29, 0, 0, 0, 0, new(time.Location))
	entry2.Entry = "Canoes are surprisingly unstable. I'm now soaked"
	entry2.Book = 3
	controller.AddEntry(entry2)
	log.Debug("entry2 added")

	for _, entry := range controller.GetPersonsEntries(user_id) {
		log.Debug(entry.String())
	}
	log.Debug("All entries retrieved")
}
