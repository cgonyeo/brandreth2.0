package importer

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"

	golog "github.com/op/go-logging"

	"github.com/dgonyeo/brandreth2.0/db"
)

var log = golog.MustGetLogger("main")

func twoDigits(num string) string {
	if len(num) < 2 {
		return "0" + num
	}
	return num
}

func sliceToString(slice []string) string {
	if len(slice) == 1 {
		return slice[0]
	}
	return slice[0] + "/" + sliceToString(slice[1:])
}

func stringToTime(timestring string) time.Time {
	tokens := strings.Split(timestring, "/")
	year, err := strconv.Atoi(tokens[0])
	if err != nil {
		log.Fatal("processing year: %v", err)
	}
	month, err := strconv.Atoi(tokens[1])
	if err != nil {
		log.Fatal("processing month: %v", err)
	}
	day, err := strconv.Atoi(tokens[2])
	if err != nil {
		log.Fatal("processing day: %v", err)
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, new(time.Location))
}

func Run(filename string, controller *db.Controller) {
	log.Debug("Reading file")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Parsing people")
	people := make(map[string](map[string]string))

	for _, line := range lines {
		people[line[0]] = make(map[string]string)
		people[line[0]]["nickname"] = line[1]
		people[line[0]]["source"] = line[5]
	}

	for name, _ := range people {
		people[name]["user_id"] = db.GetUniqueId()
	}

	log.Debug("Parsing entries")
	var entries []map[string]string

	for _, line := range lines {
		entry := make(map[string]string)
		entry["user_id"] = people[line[0]]["user_id"]
		entry["trip_reason"] = line[6]
		start := strings.Split(line[2], "/")
		entry["date_start"] = start[2] + "/" + twoDigits(start[0]) + "/" + twoDigits(start[1])
		end := strings.Split(line[3], "/")
		entry["date_end"] = end[2] + "/" + twoDigits(end[0]) + "/" + twoDigits(end[1])
		entry["entry"] = line[4]
		entry["book"] = line[7]
		entries = append(entries, entry)
	}

	log.Debug("Entering people into database")
	for name, _ := range people {
		person := db.Person{people[name]["user_id"], name, people[name]["nickname"], people[name]["source"]}
		controller.AddPerson(&person)
	}

	log.Debug("Generating trips")
	trips := make(map[string]([]map[string]string))
	for _, entry := range entries {
		placed := false
		for trip_id, entrylist := range trips {
			for _, item := range entrylist {
				if !placed && (((entry["date_start"] >= item["date_start"] && entry["date_start"] <= item["date_end"]) || (entry["date_end"] >= item["date_start"] && entry["date_end"] <= item["date_end"])) || ((item["date_start"] >= entry["date_start"] && item["date_start"] <= entry["date_end"]) || (item["date_end"] >= entry["date_start"] && item["date_end"] <= entry["date_end"]))) {
					nodup := true
					for _, item := range entrylist {
						if item["user_id"] == entry["user_id"] {
							log.Error("Error: multiple entries from the same person for a trip")
							log.Error("Name: " + controller.GetPerson(entry["user_id"]).Name)
							log.Error("Trip start: " + entry["date_start"])
							log.Error("Trip start: " + item["date_start"])
							nodup = false
						}
					}
					if nodup {
						trips[trip_id] = append(trips[trip_id], entry)
					}
					placed = true
				}
			}
		}
		if !placed {
			trip_id := db.GetUniqueId()
			trips[trip_id] = append(trips[trip_id], entry)
		}
	}

	log.Debug("Adding trips into database")
	for trip_id, entrylist := range trips {
		for _, entry := range entrylist {
			book, err := strconv.Atoi(entry["book"])
			if err != nil {
				log.Fatal("converting book: %v", err)
			}
			entrystruct := db.Entry{trip_id, entry["user_id"], entry["trip_reason"], stringToTime(entry["date_start"]), stringToTime(entry["date_end"]), entry["entry"], book}
			controller.AddEntry(&entrystruct)
		}
	}

}
