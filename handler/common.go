package handler

import (
	golog "github.com/op/go-logging"

	"github.com/dgonyeo/brandreth2.0/db"
)

var log = golog.MustGetLogger("main")

type Handler struct {
}

type PersonEntry struct {
	Person *db.Person
	Entry  *db.Entry
}
