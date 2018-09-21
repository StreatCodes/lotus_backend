package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

//TODO
type Component struct {
	ID        int
	PageID    int
	CreatedAt time.Time
	Sort      sql.NullInt64
	Data      string
}

//CreateComponentHandler returns an http handler that creates a new Component when posted to
func CreateComponentHandler(s Server) func(w http.ResponseWriter, r *http.Request) {
	type requestComponent struct {
		PageID int    `json:"page_id"`
		Sort   int    `json:"sort"`
		Data   string `json:"data"`
	}

	stmt, err := s.DB.Prepare(`INSERT INTO components (
			page_id, sort, data
		) VALUES (
			$1, $2, $3
		) RETURNING id`)

	if err != nil {
		log.Fatal("Error preparing sql statement: " + err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		var component requestComponent
		err := d.Decode(&component)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		valid := json.Valid([]byte(component.Data))
		if !valid {
			http.Error(w, `Invalid JSON found in the component data field`, http.StatusBadRequest)
			return
		}

		row := stmt.QueryRow(component.PageID, component.Sort, component.Data)

		var id int
		err = row.Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		err = enc.Encode(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
