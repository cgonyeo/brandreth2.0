package db

import (
	"fmt"
	"os"
	"strconv"

	golog "github.com/op/go-logging"
)

var log = golog.MustGetLogger("main")

func GetUniqueId() string {
	//https://groups.google.com/forum/#!topic/golang-nuts/d0nF_k4dSx4
	f, err := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	if err != nil {
		log.Info("Error opening /dev/urandom to get a unique id")
		return "lol it's broken"
	}
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uid
}

func (c *Controller) CreateTables() {
	_, err := c.getSession().Exec(setupTables)
	if err != nil {
		log.Fatal("creating tables: %v", err)
	}
}

func (c *Controller) AddPerson(person *Person) string {
	if person.UserId == "" {
		person.UserId = GetUniqueId()
	}
	_, err := c.getSession().Exec(addPerson, person.UserId, person.Name, person.Nickname, person.Source)
	if err != nil {
		log.Fatal("adding person: %v", err)
	}
	return person.UserId
}

func (c *Controller) AddEntry(entry *Entry) {
	_, err := c.getSession().Exec(addEntry, entry.TripId, entry.UserId, entry.TripReason, entry.DateStart, entry.DateEnd, entry.Entry, strconv.Itoa(entry.Book))
	if err != nil {
		log.Fatal("adding entry: %v", err)
	}
}

func (c *Controller) GetPerson(userId string) *Person {
	personData := c.getRows(getPerson, userId)
	if len(personData) == 0 {
		log.Fatal("Person not found")
	}
	if len(personData) > 1 {
		log.Fatal("Multiple people found")
	}
	person := new(Person)
	fillStruct(person, personData[0])
	return person
}

func (c *Controller) GetEntry(userId string, tripId string) *Entry {
	entryData := c.getRows(getEntry, userId, tripId)
	if len(entryData) == 0 {
		log.Fatal("Entry not found")
	}
	if len(entryData) > 1 {
		log.Fatal("Multiple entries found")
	}
	entry := new(Entry)
	fillStruct(entry, entryData[0])
	return entry
}

func (c *Controller) GetPersonsEntries(userId string) []*Entry {
	entriesData := c.getRows(getPersonsEntries, userId)
	var entries []*Entry
	for _, data := range entriesData {
		entry := new(Entry)
		fillStruct(entry, data)
		entries = append(entries, entry)
	}
	return entries
}

func (c *Controller) GetTripsEntries(tripId string) []*Entry {
	entriesData := c.getRows(getTripsEntries, tripId)
	var entries []*Entry
	for _, data := range entriesData {
		entry := new(Entry)
		fillStruct(entry, data)
		entries = append(entries, entry)
	}
	return entries
}

func (c *Controller) GetLastTrip() []*Entry {
	entryData := c.getRows(getLastEntry)
	if len(entryData) == 0 {
		log.Fatal("Entry not found")
	}
	if len(entryData) > 1 {
		log.Fatal("Multiple entries found")
	}
	entry := new(Entry)
	fillStruct(entry, entryData[0])
	return c.GetTripsEntries(entry.TripId)
}

func (c *Controller) SearchForTrips(search string) []*Entry {
	rows := c.getRows(searchQuery, c.toSearchQuery(search))
	var entries []*Entry
    for _, row := range rows {
        entries = append(entries, c.GetEntry(row["user_id"].(string), row["trip_id"].(string)))
	}
	return entries
}
