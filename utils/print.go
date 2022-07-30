package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alexandrepossebom/blocky-log/config"
	"github.com/alexandrepossebom/blocky-log/model"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jedib0t/go-pretty/v6/table"
)

func getConnection(host string) *sql.DB {
	dbconfig := config.GetDatabase(host)
	if dbconfig == nil {
		log.Println("Database not found")
		os.Exit(1)
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbconfig.Username, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.Database)
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	return db
}

func ListEventsByIP(host, ip string, hours int, eventType string) {
	db := getConnection(host)
	defer db.Close()

	var res *sql.Rows
	var err error
	if eventType == "all" {
		q := "select request_ts,client_ip,client_name,reason,response_type,question_type,question_name from log_entries where client_ip = ? and request_ts > DATE_SUB(NOW(), INTERVAL ? HOUR)"
		res, err = db.Query(q, ip, hours)
	} else {
		q := "select request_ts,client_ip,client_name,reason,response_type,question_type,question_name from log_entries where client_ip = ? and reason = ? and request_ts > DATE_SUB(NOW(), INTERVAL ? HOUR)"
		res, err = db.Query(q, ip, eventType, hours)
	}
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	entries := populateEntries(res)
	ShowTable(entries)
}

func populateEntries(res *sql.Rows) []model.LogEntry {
	entries := []model.LogEntry{}
	for res.Next() {
		var logEntry model.LogEntry
		err := res.Scan(&logEntry.RequestTime, &logEntry.ClientIP, &logEntry.ClientName, &logEntry.Reason, &logEntry.ResponseType, &logEntry.QuestionType, &logEntry.QuestionName)
		if err != nil {
			log.Panicln(err)
		}
		entries = append(entries, logEntry)
	}

	if err := res.Close(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return entries
}

func ListEvents(host, eventType string, hours int) {
	db := getConnection(host)
	defer db.Close()

	var res *sql.Rows
	var err error

	if eventType == "all" {
		res, err = db.Query("select request_ts,client_ip,client_name,reason,response_type,question_type,question_name from log_entries where request_ts > DATE_SUB(NOW(), INTERVAL ? HOUR)", hours)
	} else {
		res, err = db.Query("select request_ts,client_ip,client_name,reason,response_type,question_type,question_name from log_entries where reason = ? and request_ts > DATE_SUB(NOW(), INTERVAL ? HOUR)", eventType, hours)
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	entries := populateEntries(res)
	ShowTable(entries)
}

func ShowTable(entries []model.LogEntry) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Date", "IP", "Client Name", "Event Type", "Url"})
	m := make(map[string]time.Time)
	for _, entry := range entries {
		name := entry.ClientName
		if name == entry.ClientIP {
			name = ""
		}
		key := fmt.Sprintf("%s-%s", entry.ClientIP, entry.QuestionName)
		if _, ok := m[key]; !ok {
			m[key] = entry.RequestTime
		} else if m[key].Add(time.Second * 30).After(entry.RequestTime) {
			continue
		}
		t.AppendRow([]interface{}{entry.RequestTime.Format("02/01 15:04"), entry.ClientIP, name, entry.Reason, entry.QuestionName})
	}
	t.AppendSeparator()
	t.Render()
}
