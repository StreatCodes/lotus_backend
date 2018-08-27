package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var hasUpper = regexp.MustCompile(`[A-Z]`)
var hasLower = regexp.MustCompile(`[a-z]`)
var hasNumber = regexp.MustCompile(`[0-9]`)
var hasSymbol = regexp.MustCompile(`\W`)

func passwordValidation(password string) error {
	if len(password) < 8 {
		return errors.New("The password must be atleast 8 characters long")
	}
	if m := hasUpper.Find; m == nil {
		return errors.New("The password must contain atleast one upper case character")
	}
	if m := hasLower.Find; m == nil {
		return errors.New("The password must contain atleast one lower case character")
	}
	if m := hasNumber.Find; m == nil {
		return errors.New("The password must contain atleast one number")
	}
	if m := hasSymbol.Find; m == nil {
		return errors.New("The password must contain atleast one symbol")
	}

	return nil
}

func CreateUserHandler(s Server) func(w http.ResponseWriter, r *http.Request) {
	type requestUser struct {
		Name     string
		Email    string
		Password string
	}

	stmt, err := s.DB.Prepare(`INSERT INTO users (
			name, email, password
		) VALUES (
			$1, $2, $3
		)`)

	if err != nil {
		log.Fatal("Error preparing sql statement: " + err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		var user requestUser
		err := d.Decode(&user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = mail.ParseAddress(user.Email)
		if err != nil {
			http.Error(w, "A valid email address must be supplied", http.StatusBadRequest)
			return
		}

		err = passwordValidation(user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(user.Name) < 2 {
			http.Error(w, "The new user's name must be atleast 2 characters long", http.StatusBadRequest)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := stmt.Exec(user.Name, user.Email, passwordHash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//TODO quick hack
		fmt.Fprintf(w, `{id: %d}`, id)
	}
}
