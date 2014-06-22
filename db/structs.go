package db

import (
	"strconv"
	"time"
)

type Person struct {
    UserId   string  `sql:"user_id"`
	Name     string  `sql:"name"`
	Nickname string  `sql:"nickname"`
	Source   string  `sql:"source"`
}

type Entry struct {
    TripId     string    `sql:"trip_id"`
	UserId     string    `sql:"user_id"`
	TripReason string    `sql:"trip_reason"`
	DateStart  time.Time `sql:"date_start"`
	DateEnd    time.Time `sql:"date_end"`
	Entry      string    `sql:"entry"`
	Book       int       `sql:"book"`
}

func dateToString(date time.Time) string {
	return date.Month().String() + " " + strconv.Itoa(date.Day()) + ", " + strconv.Itoa(date.Year())
}

func (e Entry) StartString() string {
	return dateToString(e.DateStart)
}

func (e Entry) EndString() string {
	return dateToString(e.DateEnd)
}

func (p Person) String() string {
    return "UserId: " + p.UserId +"\nName: " + p.Name + "\nNickname: " + p.Nickname + "\nSource: " + p.Source
}

func (e Entry) String() string {
    return "TripId: " + e.TripId + "\nUserId: " + e.UserId + "\nTripReason: " + e.TripReason + "\nDateStart: " + e.StartString() + "\nDateEnd: " + e.EndString() + "\nEntry: " + e.Entry + "\nBook: " + strconv.Itoa(e.Book)
}
