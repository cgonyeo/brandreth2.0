package main

import (
	"time"

	golog "github.com/op/go-logging"

	"github.com/dgonyeo/brandreth2.0/db"
)

var log = golog.MustGetLogger("main")

func main() {
    test()
}

func test() {
	handler := new(db.Handler)
	log.Debug("session acquired")

	handler.CreateTables()
	log.Debug("tables created")

	person := db.Person{"", "Derek Gonyeo", "dgonz", "CSH"}
	user_id := handler.AddPerson(&person)
	log.Debug("person created")

	entry1 := new(db.Entry)
	entry1.TripId = db.GetUniqueId()
	entry1.UserId = user_id
	entry1.TripReason = "To see some awesome trees"
	entry1.DateStart = time.Date(2014, 6, 20, 0, 0, 0, 0, new(time.Location))
	entry1.DateEnd = time.Date(2014, 6, 22, 0, 0, 0, 0, new(time.Location))
	entry1.Entry = "I saw normal trees /and/ Christmas trees"
	entry1.Book = 3
	handler.AddEntry(entry1)
	log.Debug("entry1 added")

	log.Debug(handler.GetPerson(user_id).String())
	log.Debug("Person retrieved")

	log.Debug(handler.GetEntry(user_id, entry1.TripId).String())
	log.Debug("Entry retrieved")

	entry2 := new(db.Entry)
	entry2.TripId = db.GetUniqueId()
	entry2.UserId = user_id
	entry2.TripReason = "To get in a canoe"
	entry2.DateStart = time.Date(2014, 6, 27, 0, 0, 0, 0, new(time.Location))
	entry2.DateEnd = time.Date(2014, 6, 29, 0, 0, 0, 0, new(time.Location))
	entry2.Entry = "Canoes are surprisingly unstable. I'm now soaked"
	entry2.Book = 3
	handler.AddEntry(entry2)
	log.Debug("entry2 added")

    for _, entry := range handler.GetPersonsEntries(user_id) {
        log.Debug(entry.String())
    }
    log.Debug("All entries retrieved")
}
