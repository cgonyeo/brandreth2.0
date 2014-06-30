package handler

import (
	"html/template"
	"net/http"
)

type LeaderPage struct {
	Names       []string
	TripCount   []int
	Ranks       []int
	SearchQuery string
}

func (lp LeaderPage) IsActivePage(num int) bool {
	return num == 4
}

func (h *Handler) Leaderboard(w http.ResponseWriter, req *http.Request) {
	model := new(LeaderPage)

	model.Names, model.TripCount = h.c.GetLeaderboard()

	model.Ranks = make([]int, len(model.Names))

	r := 1
	lastcount := 0
	lastrank := 0
	for i, count := range model.TripCount {
		if count != lastcount {
			lastcount = count
			lastrank = r
		}
		r++
		model.Ranks[i] = lastrank
	}

	t, err := template.ParseFiles("templates/leaderboard.tmpl", "templates/stuff.tmpl")
	if err != nil {
		return
	}
	err = t.Execute(w, model)
	if err != nil {
		return
	}
}
