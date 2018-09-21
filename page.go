package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cbroglie/mustache"
)

type Page struct {
	ID        int
	CreatedAt time.Time
	Title     string
	Slug      string
	Parent    sql.NullInt64
	Sort      sql.NullInt64
}

//Recursively loop over pages to build the full page URL
func (p *Page) buildURL(pages []Page) string {
	if p.Parent.Valid {
		var parent Page
		//Find the parent page in pages
		for _, v := range pages {
			if v.ID == int(p.Parent.Int64) {
				parent = v
			}
		}
		return parent.buildURL(pages) + "/" + p.Slug
	}
	return "/" + p.Slug
}

func (p *Page) buildHTML(s Server) (string, error) {
	rows, err := s.DB.Query(
		`SELECT data FROM components WHERE page_id=$1 ORDER BY sort`,
		p.ID)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	html := ""
	for rows.Next() {
		var pageData []byte

		err := rows.Scan(&pageData)
		if err != nil {
			log.Fatal(err)
		}

		var data interface{}
		err = json.Unmarshal(pageData, &data)
		if err != nil {
			return "", err
		}

		res, err := mustache.RenderFile("./templates/sample.html", data)
		if err != nil {
			return "", err
		}

		html += res
	}

	fmt.Println(html)
	return html, nil
}

func buildPages(s Server) error {
	fmt.Println("Running full page build")
	rows, err := s.DB.Query(`SELECT id, title, slug, parent FROM pages`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var pages []Page

	for rows.Next() {
		var page Page

		err := rows.Scan(&page.ID, &page.Title, &page.Slug, &page.Parent)
		if err != nil {
			log.Fatal(err)
		}
		pages = append(pages, page)
	}

	siteTree := make(map[string]Page)

	for _, page := range pages {
		url := page.buildURL(pages)
		siteTree[url] = page

		fmt.Printf("Building: %s\n", url)
		page.buildHTML(s)
	}

	return nil
}

//ServePageHandler serves the webpages from memory to a user
func ServePageHandler(s Server) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		//Get list of encodings client supports, if they support gzip
		//then serve them gzipped content
		encodings := strings.ToLower(r.Header.Get("Accept-Encoding"))

		supportsGzip := strings.Contains(encodings, "gzip")
		if supportsGzip {
			fmt.Fprint(w, r.URL.RequestURI())
			fmt.Fprint(w, "you support gzip")
		} else {
			fmt.Fprint(w, "you don't support gzip")
		}
	}
}

//CreatePageHandler returns an http handler that creates a new Page when posted to
func CreatePageHandler(s Server) func(w http.ResponseWriter, r *http.Request) {
	type requestPage struct {
		Title  string        `json:"title"`
		Slug   string        `json:"slug"`
		Parent sql.NullInt64 `json:"parent"`
		Sort   int           `json:"sort"`
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
