package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func authorizeHandler(s Server) func(w http.ResponseWriter, r *http.Request) {
	type login struct {
		email    string
		password string
	}

	stmt, err := s.DB.Prepare(`SELECT id, password FROM users WHERE email=$1`)
	if err != nil {
		log.Fatal("Error preparing sql statement: " + err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		var requestLogin login
		err := d.Decode(&requestLogin)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if requestLogin.email == "" || requestLogin.password == "" {
			http.Error(w, "An email and password must be present in the request body.", http.StatusBadRequest)
			return
		}

		var id int
		var hash string
		row := stmt.QueryRow(requestLogin.email)

		err = row.Scan(&id, &hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(requestLogin.password))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		key := make([]byte, 40)
		_, err = rand.Read(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%010d:%#p", id, key)
	}
}
