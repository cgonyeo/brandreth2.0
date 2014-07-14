package handler

import (
    "encoding/json"
    "time"
	"html/template"
	"net/http"

    "github.com/dgonyeo/brandreth2.0/db"
)

type NewTripPage struct {
	Names       []string
	Reasons     []string
	SearchQuery string
}

func (ntp NewTripPage) IsActivePage(num int) bool {
	return false
}

func (h *Handler) NewTrip(w http.ResponseWriter, req *http.Request) {
	model := new(NewTripPage)

	people := h.c.GetPeople()

	model.Names = make([]string, len(people))
	model.Reasons = h.c.GetTripReasons()

	for i, person := range people {
		model.Names[i] = person.Name
	}

	t, err := template.ParseFiles("templates/newtrip.tmpl", "templates/stuff.tmpl")
	if err != nil {
		return
	}
	err = t.Execute(w, model)
	if err != nil {
		return
	}
}

var tempTrip *TripJson

type TripEntryJson struct {
	Name      string `json:"name"`
    Book int `json:"book"`
	Arrival   string `json:"arrival"`
	Departure string `json:"departure"`
	Entry     string `json:"entry"`
}

type TripJson struct {
    Reason string `json:"reason"`
    Entries []TripEntryJson `json:"entries"`
}

func (h *Handler) SubmitTrip(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    var trip TripJson
    err := decoder.Decode(&trip)
    if err != nil {
        log.Error("Couldn't decode json: %v", err)
        return
    }
    reason := trip.Reason
    tripId := db.GetUniqueId()
    for _, entryData := range trip.Entries {
        dateStart, err := time.Parse("2006-01-02", entryData.Arrival)
        if err != nil {
            log.Error("Couldn't parse date: %v", err)
            return
        }

        dateEnd, err := time.Parse("2006-01-02", entryData.Departure)
        if err != nil {
            log.Error("Couldn't parse date: %v", err)
            return
        }

        userId := h.c.GetUserIdByName(entryData.Name)
        if userId == "" {
            log.Debug("Brandreth noob detection algorithm triggered")
            w.WriteHeader(http.StatusOK)
            //TODO: search through list for rest of noobs, write in json
            return
        }

        entry := &db.Entry{tripId, userId, reason, dateStart, dateEnd, entryData.Entry, entryData.Book}
        h.c.AddEntry(entry)
        log.Debug("Adding: \n%s", entry.String())
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Success"))
}
