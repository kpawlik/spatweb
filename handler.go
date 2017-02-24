package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

func getRecords(rows *sql.Rows, w http.ResponseWriter) results {
	var count int

	cols, err := rows.Columns()
	if err != nil {
		fmt.Fprintf(w, "Get colmns error: %v\n", err)
		log.Printf("Get colmns error:  %v", err)
	}
	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}
	recs := recs{}
	for rows.Next() {
		count++
		rec := make(jsonRec)
		err = rows.Scan(vals...)
		for i, val := range vals {
			switch val.(type) {
			case *sql.RawBytes, sql.RawBytes:
				s := val.(*sql.RawBytes)
				ss := []byte(*s)
				if cols[i] != "the_geom" {
					rec[cols[i]] = string(ss)
				}
			}
		}
		recs = append(recs, rec)
	}
	res := make(results)
	res["records"] = recs
	res["count"] = count
	return res
}

type handler struct {
	db *sql.DB
}

func newHandler(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) empty(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) getRecordsByValue(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)
	vars := mux.Vars(r)
	table := vars["table"]
	column := vars["column"]
	value := vars["value"]

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error (getRecordsByValue): %v\n", err)
		}
	}()

	where := fmt.Sprintf(" %s = '%s'", column, value)
	query := fmt.Sprintf(baseQuery+"  WHERE %s", table, where)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Fprintf(w, "Run query error: %v\n", err)
		log.Printf("Query error '%s'; %v", query, err)
		return
	}
	defer rows.Close()
	recs := getRecords(rows, w)
	jsonEnc := json.NewEncoder(w)

	jsonEnc.Encode(recs)
}

func (h *handler) getRecords(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error (getRecords): %v\n", err)
		}
	}()
	vars := mux.Vars(r)
	table := vars["table"]

	query := fmt.Sprintf(baseQuery, table)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Fprintf(w, "Run query error: %v\n", err)
		log.Printf("Query error '%s'; %v", query, err)
		return
	}
	defer rows.Close()
	recs := getRecords(rows, w)
	jsonEnc := json.NewEncoder(w)

	jsonEnc.Encode(recs)
}
func (h *handler) getRecordsLimit(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error (getRecordsLimit): %v\n", err)
		}
	}()
	vars := mux.Vars(r)
	table := vars["table"]
	limit := vars["limit"]

	query := fmt.Sprintf(baseQuery+"LIMIT %s", table, limit)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Fprintf(w, "Run query error: %v\n", err)
		log.Printf("Query error '%s'; %v", query, err)
		return
	}
	defer rows.Close()
	recs := getRecords(rows, w)
	jsonEnc := json.NewEncoder(w)

	jsonEnc.Encode(recs)
}
