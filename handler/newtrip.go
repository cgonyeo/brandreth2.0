package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/dgonyeo/brandreth2.0/config"
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

func IsNotAnAdmin(req *http.Request) bool {
	username, ok := req.Header["X-Webauth-User"]
	if !ok {
		return true
	}
	for _, name := range config.Config.Admins.Name {
		if username[0] == name {
			return false
		}
	}
	return true
}

func (h *Handler) NewTrip(w http.ResponseWriter, req *http.Request) {
	if IsNotAnAdmin(req) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	model := new(NewTripPage)

	people := h.c.GetPeople()

	model.Names = make([]string, len(people))
	model.Reasons = h.c.GetTripReasons()

	for i, person := range people {
		model.Names[i] = person.Name
	}

	t, err := template.ParseFiles(
		config.Config.Templates.Path+"templates/newtrip.tmpl",
		config.Config.Templates.Path+"templates/stuff.tmpl")
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
	Book      int    `json:"book"`
	Arrival   string `json:"arrival"`
	Departure string `json:"departure"`
	Entry     string `json:"entry"`
}

type TripJson struct {
	Reason  string          `json:"reason"`
	Entries []TripEntryJson `json:"entries"`
}

type NubzReturnJson struct {
	Nubz   []TripEntryJson `json:"noobs"`
	TripId string          `json:"trip_id"`
}

func (h *Handler) SubmitTrip(w http.ResponseWriter, req *http.Request) {
	if IsNotAnAdmin(req) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(req.Body)
	var trip TripJson
	err := decoder.Decode(&trip)
	if err != nil {
		log.Error("Couldn't decode json: %v", err)
		return
	}
	var nubz []TripEntryJson
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
			nubz = append(nubz, entryData)
		} else {
			entry := &db.Entry{tripId, userId, reason, dateStart, dateEnd, entryData.Entry, entryData.Book}
			h.c.AddEntry(entry)
			log.Debug("Adding: \n%s", entry.String())
		}
	}
	w.WriteHeader(http.StatusOK)
	if len(nubz) == 0 {
		w.Write([]byte("{\"trip_id\": \"" + tripId + "\", \"success\": \"t\"}"))
	} else {
		n := &NubzReturnJson{nubz, tripId}
		content, err := json.Marshal(n)
		if err != nil {
			log.Error("Marshaling nubz!: %v", err)
			return
		}
		w.Write(content)
	}
}

type NubJson struct {
	Name      string `json:"name"`
	Book      int    `json:"book"`
	Arrival   string `json:"arrival"`
	Departure string `json:"departure"`
	Entry     string `json:"entry"`
	Nickname  string `json:"nickname"`
	Source    string `json:"source"`
}

type NubzJson struct {
	TripId string    `json:"trip_id"`
	Nubz   []NubJson `json:"nubz"`
}

func (h *Handler) SubmitNubz(w http.ResponseWriter, req *http.Request) {
	if IsNotAnAdmin(req) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(req.Body)
	var nubz NubzJson
	err := decoder.Decode(&nubz)
	if err != nil {
		log.Error("Couldn't decode json: %v", err)
		return
	}
	for _, nub := range nubz.Nubz {
		dateStart, err := time.Parse("2006-01-02", nub.Arrival)
		if err != nil {
			log.Error("Couldn't parse date: %v", err)
			return
		}

		dateEnd, err := time.Parse("2006-01-02", nub.Departure)
		if err != nil {
			log.Error("Couldn't parse date: %v", err)
			return
		}

		userId := h.c.AddPerson(&db.Person{
			UserId:   db.GetUniqueId(),
			Name:     nub.Name,
			Nickname: nub.Nickname,
			Source:   nub.Source,
		})

		entry := &db.Entry{nubz.TripId, userId, h.c.GetTripReason(nubz.TripId), dateStart, dateEnd, nub.Entry, nub.Book}
		h.c.AddEntry(entry)
		log.Debug("Adding: \n%s", entry.String())
	}
	w.Write([]byte("{\"trip_id\": \"" + nubz.TripId + "\", \"success\": \"t\"}"))
}
