package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Init() {
	a.initDBConn()
	a.setupTables()
	a.appendRoutes()
}

func (a *App) appendRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", a.handleIndex).Methods("GET")
	r.HandleFunc("/{category}", a.handleCategory).Methods("GET")
	a.Router = r
}

func (a *App) initDBConn() {
	host := os.Getenv("pg_host")
	port := os.Getenv("pg_port")
	convPort, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		fmt.Println("Issue with PORT ENVVAR")
	}
	user := os.Getenv("pg_user")
	password := os.Getenv("pg_password")
	dbname := os.Getenv("pg_dbname")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, convPort, user, password, dbname)

	var dbErr error
	a.DB, dbErr = sql.Open("postgres", psqlInfo)
	if dbErr != nil {
		fmt.Println("DB Conn error")
	}
	fmt.Println("DB Conn Success")
}

func (a *App) setupTables() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS employee
	(
		id SERIAL,
		name TEXT NOT NULL,
		department TEXT NOT NULL,
		title TEXT NOT NULL,
		remuneration NUMERIC(10, 2),
		expenses NUMERIC(10, 2),
		year INT,
		CONSTRAINT employeee_pkey PRIMARY KEY (id)
	)`
	if _, err := a.DB.Exec(createTableQuery); err != nil {
		//"274,262,559","1,284,768"
		fmt.Println("Error with creating 'remuneration' table.")
		log.Fatal(err)
	}
	//a.populateTables()
}

func (a *App) populateTables() {
	path := "./2017StaffRemunerationOver75KWithExpenses.csv"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(bufio.NewReader(file))

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("file read err")
		}

		salary, salErr := strconv.ParseFloat(strings.Replace(line[3], ",", "", -1), 64)
		if salErr != nil {
			fmt.Println(err)
		}
		expense, expErr := strconv.ParseFloat(strings.Replace(line[4], ",", "", -1), 64)
		if expErr != nil {
			fmt.Println(err)
		}
		latest := Employee{
			Name:         line[0],
			Department:   line[1],
			Title:        line[2],
			Remuneration: salary,
			Expenses:     expense,
			Year:         2018,
		}
		if err := latest.createEmployee(a.DB); err != nil {
			fmt.Println("Issue with populating DB.")
			fmt.Println(err)
		}
	}
}

func (a *App) Start() {
	fmt.Println("Listening on PORT 9123")
	log.Fatal(http.ListenAndServe(":9123", a.Router))
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	val, err := getAllEmployees(a.DB)
	if err != nil {
		fmt.Println(err)
	}
	respondWithJSON(w, http.StatusOK, val)
}

func (a *App) handleCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queryFilter := vars["category"]
	val, err := getSomeEmployees(a.DB, queryFilter)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, val)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, message)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	type message struct {
		Date  string      `json:"date"`
		Data  interface{} `json:"data"`
		Notes string      `json:"notes"`
	}
	res, err := json.Marshal(
		message{time.Now().Format("02-01-2006"),
			payload,
			"No notes added.",
		})
	if err != nil {
		fmt.Println("Error with Responding in JSON and marshalling of data.")
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
}
