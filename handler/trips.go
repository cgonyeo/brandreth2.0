package handler

import (
	"html/template"
	"net/http"

	"github.com/mholt/binding"

	"github.com/dgonyeo/brandreth2.0/config"
)

type TripsPage struct {
	TripItems    []*TripPageItem
	SearchQuery  string
	TripsPageNum int
	PagesToShow  []int
	NumPages     int
}

func (tp TripsPage) IsActivePage(num int) bool {
	return num == 1
}

func (tp TripsPage) IsTripsPageNum(num int) bool {
	return num == tp.TripsPageNum
}

func (tp TripsPage) IsLastPage() bool {
	return tp.TripsPageNum == tp.NumPages
}

func (tp TripsPage) PrevPage() int {
	return tp.TripsPageNum - 1
}

func (tp TripsPage) NextPage() int {
	return tp.TripsPageNum + 1
}

type TripPageItem struct {
	Start       string
	End         string
	Reason      string
	NumPeople   int
	TripId      string
	SearchQuery string
}

type TripsParams struct {
	Page int
}

func (tp *TripsParams) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&tp.Page: binding.Field{
			Form:     "page",
			Required: true,
		},
	}
}

func (h *Handler) Trips(w http.ResponseWriter, req *http.Request) {
	tripsParams := new(TripsParams)
	errs := binding.Bind(req, tripsParams)
	if errs.Handle(w) {
		log.Error("Error with binding")
		return
	}

	page := tripsParams.Page

	trips := h.c.GetRecentTrips(10, page)

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

	tp.TripsPageNum = page
	tp.NumPages = h.c.GetNumPages(10)
	for i := page - 5; i < page+5; i++ {
		if i >= 0 && i <= tp.NumPages {
			tp.PagesToShow = append(tp.PagesToShow, i)
		}
	}

	t, err := template.ParseFiles(
		config.Config.Templates.Path+"templates/trips.tmpl",
		config.Config.Templates.Path+"templates/stuff.tmpl")
	if err != nil {
		return
	}
	err = t.Execute(w, tp)
	if err != nil {
		return
	}
}
