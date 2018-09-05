package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

//CreatePageHandler returns an http handler that creates a new Page when posted to
func CreatePageHandler(s Server) func(w http.ResponseWriter, r *http.Request) {
	type requestPage struct {
		Title  string `json:"title"`
		Slug   string `json:"slug"`
		Parent int    `json:"parent"`
		Sort   int    `json:"sort"`
	}

	stmt, err := s.DB.Prepare(`INSERT INTO pages (
			title, slug, parent, sort
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id`)

	if err != nil {
		log.Fatal("Error preparing sql statement: " + err.Error())
	}

	slugMatch, err := regexp.Compile(`[A-Za-z0-9\-._]{1,40}`)
	if err != nil {
		log.Fatal("Error compiling regexp: " + err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		var page requestPage
		err := d.Decode(&page)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		match := slugMatch.MatchString(page.Slug)
		if !match {
			http.Error(w, `The page slug must be between 1 and 40 characters in length\
				and only contain characters in the following set: A-Z a-z 0-9 - . _`, http.StatusBadRequest)
			return
		}

		row := stmt.QueryRow(page.Title, page.Slug, page.Parent, page.Sort)

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
