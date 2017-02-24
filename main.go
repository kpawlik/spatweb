package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"runtime"

	"github.com/gorilla/mux"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type jsonRec map[string]string
type recs []jsonRec
type results map[string]interface{}

var (
	db        *sql.DB
	dbPath    string
	logPath   string
	dbDriver  string
	baseQuery string
)

func init() {
	sql.Register("spatialite", &sqlite3.SQLiteDriver{
		Extensions: []string{
			"libspatialite",
		},
	})
	if runtime.GOOS == "windows" {
		dbPath = "c:\\tmp\\atlanta_1601.db"
		logPath = "c:\\tmp\\sqliteapp.log"
		dbDriver = "sqlite3"
		baseQuery = "select * from data$%s "
	} else {
		dbPath = "/kpa-tmp/atlanta_1601.db"
		logPath = "/kpa-tmp/sqliteapp.log"
		dbDriver = "spatialite"
		baseQuery = "select *, st_astext(the_geom) AS geom from data$%s "
	}

}

func main() {
	var (
		err     error
		logFile *os.File
	)
	// setup log
	logFile, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("server started")
	// setup db
	if db, err = sql.Open(dbDriver, dbPath); err != nil {
		log.Printf("Driver error: %v\n", err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Printf("Ping  error: %v\n", err)
	}

	router := mux.NewRouter()

	handler := newHandler(db)
	router.HandleFunc("/favicon.ico", handler.empty)
	router.HandleFunc("/{table}/{column:\\w+}={value:.*}", handler.getRecordsByValue)
	router.HandleFunc("/{table}", handler.getRecords)
	router.HandleFunc("/{table}/{limit:\\d+}", handler.getRecordsLimit)

	http.ListenAndServe(":8080", router)
}
