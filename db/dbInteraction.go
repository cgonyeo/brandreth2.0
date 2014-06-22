package db

import (
	"database/sql"
	_ "github.com/lib/pq"
    "reflect"
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

var addPerson = "INSERT INTO people (user_id,name,nickname,source) VALUES ($1,$2,$3,$4);"

var addEntry = "INSERT INTO entries (trip_id,user_id,trip_reason,date_start,date_end,entry,book) VALUES ($1,$2,$3,$4,$5,$6,$7);"

var getEntry = "SELECT * FROM entries WHERE user_id=$1 AND trip_id=$2;"

var getPerson = "SELECT * FROM people WHERE user_id=$1;"

var getPersonsEntries = "SELECT * FROM entries WHERE user_id=$1;"

var getTripsEntries = "SELECT * FROM entries WHERE trip_id=$1;"

var getAllEntries = "SELECT * FROM entries;"

type Handler struct {
	db *sql.DB
}

func (h Handler) getSession() *sql.DB {
	if h.db == nil {
		db, err := sql.Open("postgres", "postgres://brandreth:password@localhost/brandreth?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
		h.db = db
	}
	return h.db
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

func (h Handler) getRows(queryString string, args ...interface{}) []map[string]interface{} {
	rows, err := h.getSession().Query(queryString, args...)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
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
