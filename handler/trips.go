package handler

import (
	"html/template"
	"net/http"
)

type TripsPage struct {
	TripItems   []*TripPageItem
	SearchQuery string
}

func (tp TripsPage) IsActivePage(num int) bool {
	log.Debug("TripsPage")
	return num == 1
}

type TripPageItem struct {
	Start       string
	End         string
	Reason      string
	NumPeople   int
	TripId      string
	SearchQuery string
}

func (h *Handler) Trips(w http.ResponseWriter, req *http.Request) {
	trips := h.c.GetRecentTrips(10, 1)

	tp := new(TripsPage)
	for _, trip := range trips {
		earliestStart := trip[0]
		latestEnd := trip[0]
		tripReason := trip[0].TripReason
		tripId := trip[0].TripId
		for _, entry := range trip {
			if entry.DateStart.Before(earliestStart.DateStart) {
				earliestStart = entry
			}
			if entry.DateEnd.After(latestEnd.DateEnd) {
				latestEnd = entry
			}
		}

		tripItem := new(TripPageItem)
		tripItem.Start = earliestStart.StartString()
		tripItem.End = earliestStart.EndString()
		tripItem.Reason = tripReason
		tripItem.NumPeople = len(trip)
		tripItem.TripId = tripId

		tp.TripItems = append(tp.TripItems, tripItem)
	}

	t, err := template.ParseFiles("templates/trips.tmpl", "templates/stuff.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, tp)
	if err != nil {
		log.Fatal(err)
	}
}
