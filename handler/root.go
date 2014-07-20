package handler

import (
	"html/template"
	"net/http"

	"github.com/dgonyeo/brandreth2.0/config"
)

type RootPage struct {
	SearchQuery string
}

func (pp RootPage) IsActivePage(num int) bool {
	log.Debug("root")
	return false
}

func (h *Handler) Root(w http.ResponseWriter, req *http.Request) {
	model := new(RootPage)

	t, err := template.ParseFiles(
		config.Config.Templates.Path+"templates/root.tmpl",
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
