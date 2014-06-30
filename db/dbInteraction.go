package db

import (
	"database/sql"
	"reflect"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
)

var setupTables = `
BEGIN;
    drop table if exists entries;
    drop table if exists people;

    create table people (
        user_id TEXT NOT NULL,
        name TEXT NOT NULL,
        nickname TEXT NOT NULL,
        source TEXT,
        PRIMARY KEY (user_id)
    )
    WITHOUT OIDS;

    create table entries (
        trip_id TEXT NOT NULL,
        user_id TEXT NOT NULL,
        trip_reason TEXT,
        date_start DATE NOT NULL,
        date_end DATE NOT NULL,
        entry TEXT NOT NULL,
        book INTEGER NOT NULL,
        PRIMARY KEY (trip_id, user_id),
        CONSTRAINT user_id_fkey FOREIGN KEY (user_id)
            REFERENCES people (user_id) MATCH SIMPLE
            ON UPDATE NO ACTION ON DELETE NO ACTION
    )
    WITHOUT OIDS;
COMMIT;
`

var searchQuery = `
SELECT trip_id, user_id
FROM (SELECT entries.trip_id as trip_id,
        people.user_id as user_id,
        to_tsvector(people.name) || 
        to_tsvector(people.nickname) || 
        to_tsvector(entries.trip_reason) || 
        to_tsvector(entries.entry) as document 
    FROM entries 
    JOIN people ON people.user_id = entries.user_id 
    GROUP BY 
        people.user_id,
        people.name, 
        people.nickname, 
        entries.trip_id,
        entries.trip_reason,
        entries.entry,
        entries.date_end
        ORDER BY entries.date_end DESC) e_search 
WHERE e_search.document @@ to_tsquery($1);
`

var addPerson = "INSERT INTO people (user_id,name,nickname,source) VALUES ($1,$2,$3,$4);"

var addEntry = "INSERT INTO entries (trip_id,user_id,trip_reason,date_start,date_end,entry,book) VALUES ($1,$2,$3,$4,$5,$6,$7);"

var getEntry = "SELECT * FROM entries WHERE user_id=$1 AND trip_id=$2;"

var getPerson = "SELECT * FROM people WHERE user_id=$1;"

var getPeople = "SELECT * FROM people ORDER BY people.name ASC;"

var getPersonsEntries = "SELECT * FROM entries WHERE user_id=$1 ORDER BY date_end DESC;"

var getTripIds = "SELECT DISTINCT trip_id, date_end FROM entries ORDER BY date_end DESC LIMIT $1 OFFSET $2;"

var getTripsEntries = "SELECT * FROM entries WHERE trip_id=$1 ORDER BY date_end DESC;"

var getAllEntries = "SELECT * FROM entries ORDER BY date_end DESC;"

var getLastEntry = "SELECT * FROM entries ORDER BY date_end DESC LIMIT 1;"

var getPeoplesFirstTrip = " SELECT count(user_id), date_part FROM (SELECT user_id, date_part('year', min(date_start)) FROM entries GROUP BY user_id) AS q GROUP BY date_part ORDER BY date_part;"

var getPeoplePerYear = "SELECT count(*), date_part('year', date_start) FROM entries GROUP BY date_part ORDER BY date_part;"

var getUniquePeoplePerYear = "SELECT count(DISTINCT user_id), date_part('year', date_start) FROM entries GROUP BY date_part ORDER BY date_part;"

var getDaysAtBrandrethPerYear = "SELECT date_part('year', date_start), sum(date_end - date_start) AS duration from entries GROUP BY date_part ORDER BY date_part;"

var numFromSourcesPerYear = "select count(user_id), source, date_part('year', date_start) from people INNER JOIN entries USING (user_id) GROUP BY source, date_part ORDER BY date_part;"

var avgDurationPerYear = "SELECT date_part('year', date_start), avg(date_end - date_start) AS duration from entries GROUP BY date_part ORDER BY date_part;"

var getSourceBreakdown = "SELECT count(user_id), source FROM people GROUP BY source ORDER BY source;"

var getMonthsForAPerson = "select date_part('month', date_start), count(trip_id) from entries where user_id=$1 GROUP BY date_part ORDER BY date_part;"

var getTripCountForEveryone = "SELECT count(trip_id), name from entries INNER JOIN people USING (user_id) GROUP BY name ORDER BY count DESC;"

var getReasons = "SELECT DISTINCT trip_reason FROM entries ORDER BY trip_reason;"

var getUserByName = "SELECT user_id FROM PEOPLE WHERE name=$1;"

type Controller struct {
	db *sql.DB
}

func (c *Controller) toSearchQuery(search string) string {
	re1, _ := regexp.Compile("[\\(\\)\\&\\|\\!'\"]+")
	re2, _ := regexp.Compile("[ 	\n]+")
	return re2.ReplaceAllString(re1.ReplaceAllString(strings.Trim(search, " 	\n"), ""), " & ")
}

func (c *Controller) getSession() *sql.DB {
	if c.db == nil {
		log.Debug("opening connection")
		db, err := sql.Open("postgres", "postgres://brandreth:password@localhost/brandreth?sslmode=disable")
		if err != nil {
			log.Fatal("actuiring session: %v", err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatal("pinging after acuiring session: %v", err)
		}
		c.db = db
	}
	return c.db
}

func fillStruct(toFill interface{}, data map[string]interface{}) {
	ftd := make(map[string]reflect.Value)
	typeOf := reflect.TypeOf(toFill).Elem()
	valueOf := reflect.ValueOf(toFill)
	for i := 0; i < typeOf.NumField(); i++ {
		ftd[typeOf.Field(i).Tag.Get("sql")] = reflect.Indirect(valueOf).Field(i)
	}

	for key, val := range data {
		ftd[key].Set(reflect.ValueOf(val).Convert(ftd[key].Type()))
	}
}

func (c *Controller) getRows(queryString string, args ...interface{}) []map[string]interface{} {
	rows, err := c.getSession().Query(queryString, args...)
	defer rows.Close()

	if err != nil {
		log.Fatal("Querying rows: %v", err)
	}

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	var returnSlice []map[string]interface{}

	//To get all values we'd do `for rows.Next {`, but I only care about the first row
	for rows.Next() {
		returnMap := make(map[string]interface{})
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			returnMap[col] = v
		}
		returnSlice = append(returnSlice, returnMap)
	}
	return returnSlice

}
