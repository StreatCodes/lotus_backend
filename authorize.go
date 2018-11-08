package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//UserContextKey namespace our context from other libraries
type UserContextKey string

func authorizeHandler(s Server) func(w http.ResponseWriter, r *http.Request) {
	type login struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"remember_me"`
	}

	selectUserStmt, err := s.DB.Prepare(`SELECT id, password FROM users WHERE email=$1`)
	if err != nil {
		log.Fatal("Error preparing sql statement: " + err.Error())
	}

	createSessionStmt, err := s.DB.Prepare(`INSERT INTO sessions (user_id, created_ip, key, expires) VALUES ($1, $2, $3, $4)`)
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

		if requestLogin.Email == "" || requestLogin.Password == "" {
			http.Error(w, "An email and password must be present in the request body.", http.StatusBadRequest)
			return
		}

		//Get user information and password hash from DB
		var id int
		var hash string
		row := selectUserStmt.QueryRow(requestLogin.Email)

		err = row.Scan(&id, &hash)
		if err == sql.ErrNoRows {
			http.Error(w, "No account found with the given email address.", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Compare password with bcrypt hash
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(requestLogin.Password))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		//Create key based on 0 padded user id (ten in length) : 40 bytes as hex
		randBytes := make([]byte, 40)
		_, err = rand.Read(randBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key := fmt.Sprintf("%x", randBytes)

		//If remember me is set keep session for a year, otherwise 2 hours
		var expires time.Time
		if requestLogin.RememberMe {
			expires = time.Now().Add(time.Hour * 24 * 365)
		} else {
			expires = time.Now().Add(time.Hour * 2)
		}

		_, err = createSessionStmt.Exec(id, r.RemoteAddr, key, expires)
		if err != nil {
			http.Error(w, "Error creating session"+err.Error(), http.StatusInternalServerError)
		}

		//Write key to response with the user_id included
		enc := json.NewEncoder(w)
		err = enc.Encode(fmt.Sprintf("%010d:%s", id, key))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func authorizedMiddleware(s Server) func(next http.Handler) http.Handler {
	selectUserStmt, err := s.DB.Prepare(`SELECT user_id FROM sessions WHERE user_id=$1 AND key=$2 AND expires>NOW()`)
	if err != nil {
		log.Fatal("Error preparing sql statement: " + err.Error())
	}

	//A strict regex to get the authorization token from the header
	authRegex, err := regexp.Compile(`LotusToken: (\d{10}):([0-9a-f]{80}$)`)
	if err != nil {
		log.Fatal("Error compiling auth token regexp:" + err.Error())
	}

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authString := r.Header.Get("Authorization")
			matches := authRegex.FindStringSubmatch(authString)

			if len(matches) != 3 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			id := strings.TrimLeft(matches[1], "0")
			key := matches[2]

			row := selectUserStmt.QueryRow(id, key)

			var userID int
			err := row.Scan(&userID)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			context.WithValue(r.Context(), UserContextKey("userID"), userID)

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
