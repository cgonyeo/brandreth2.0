package handler

import (
	"html/template"
	"net/http"

	"github.com/dgonyeo/brandreth2.0/config"
)

type StatsPage struct {
	Years                        []int
	NumPeoples                   []int
	NumNoobs                     []int
	NumUniqueVisitors            []int
	NumDays                      []int
	Sources                      []string
	NumVisitorsForSourcesPerYear [][]int
	Durations                    []float64
	SourcesForPie                []string
	PeopleForPie                 []int
	SearchQuery                  string
}

func (pp StatsPage) IsActivePage(num int) bool {
	return num == 3
}

func (h *Handler) Stats(w http.ResponseWriter, req *http.Request) {
	model := new(StatsPage)

	model.Years, model.NumPeoples = h.c.GetYearsToNumVisitors()
	_, model.NumNoobs = h.c.GetYearsToNumNewVisitors()
	_, model.NumUniqueVisitors = h.c.GetYearsToUniqueVisitors()
	_, model.NumDays = h.c.GetYearsToDays()
	_, model.Sources, model.NumVisitorsForSourcesPerYear = h.c.GetYearsToVisitorsSources()
	_, model.Durations = h.c.GetAvgDurationPerYear()
	model.SourcesForPie, model.PeopleForPie = h.c.GetSources()

	t, err := template.ParseFiles(
		config.Config.Templates.Path+"templates/stats.tmpl",
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
