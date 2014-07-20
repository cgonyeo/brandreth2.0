package handler

import (
	"html/template"
	"net/http"

	"github.com/dgonyeo/brandreth2.0/config"
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

	t, err := template.ParseFiles(
		config.Config.Templates.Path+"templates/leaderboard.tmpl",
		config.Config.Templates.Path+"templates/stuff.tmpl")
	if err != nil {
		log.Error("Error parsing the templates: %v", err)
		return
	}
	err = t.Execute(w, model)
	if err != nil {
		log.Error("Error executing the templates: %v", err)
		return
	}
}
