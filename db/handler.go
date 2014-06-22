package db

import (
	"fmt"
	golog "github.com/op/go-logging"
	"os"
	"strconv"
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

func (h Handler) CreateTables() {
	_, err := h.getSession().Exec(setupTables)
	if err != nil {
		log.Fatal(err)
	}
}

func (h Handler) AddPerson(person *Person) string {
	if person.UserId == "" {
		person.UserId = GetUniqueId()
	}
	_, err := h.getSession().Exec(addPerson, person.UserId, person.Name, person.Nickname, person.Source)
	if err != nil {
		log.Fatal(err)
	}
	return person.UserId
}

func (h Handler) AddEntry(entry *Entry) {
	_, err := h.getSession().Exec(addEntry, entry.TripId, entry.UserId, entry.TripReason, entry.DateStart, entry.DateEnd, entry.Entry, strconv.Itoa(entry.Book))
	if err != nil {
		log.Fatal(err)
	}
}

func (h Handler) GetPerson(userId string) *Person {
    //personData := h.getRowContents("SELECT * FROM people WHERE user_id='" + userId + "';")
    personData := h.getRows(getPerson, userId)
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

func (h Handler) GetEntry(userId string, tripId string) *Entry {
    entryData := h.getRows(getEntry, userId, tripId)
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

func (h Handler) GetPersonsEntries(userId string) []*Entry {
    entriesData := h.getRows(getPersonsEntries, userId)
    var entries []*Entry
    for _, data := range entriesData {
        entry := new(Entry)
        fillStruct(entry, data)
        entries = append(entries, entry)
    }
    return entries
}

func (h Handler) GetTripsEntries(tripId string) []*Entry {
    entriesData := h.getRows(getTripsEntries, tripId)
    var entries []*Entry
    for _, data := range entriesData {
        entry := new(Entry)
        fillStruct(entry, data)
        entries = append(entries, entry)
    }
    return entries
}

func (h Handler) GetAllEntries() []*Entry {
    entriesData := h.getRows(getAllEntries)
    var entries []*Entry
    for _, data := range entriesData {
        entry := new(Entry)
        fillStruct(entry, data)
        entries = append(entries, entry)
    }
    return entries
}
