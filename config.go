package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/mail"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

//Server is our global struct which contains configuration information
type Server struct {
	Env      string
	HTTPAddr string
	HTTPPort string
	DBUser   string
	DBPass   string
	DBAddr   string
	DBName   string
	DB       *sql.DB
	SiteTree SiteTree
}

//PostgresURL is a convience method that returns the postgres connection URL
func (s *Server) PostgresURL() string {
	return "postgres://" +
		s.DBUser + ":" +
		s.DBPass + "@" +
		s.DBAddr + "/" +
		s.DBName + "?sslmode=disable"
}

//HTTPAddress is a convience method that returns the HTTP address that the web server should bind to
func (s *Server) HTTPAddress() string {
	return s.HTTPAddr + ":" + s.HTTPPort
}

func generatePassword() string {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()"

	b := make([]byte, 10)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}

	return string(b)
}

func firstTimeSetup(db *sql.DB) error {
	var email, confirmEmail string

	fmt.Print("Enter Email: ")
	_, err := fmt.Scanln(&email)
	if err != nil {
		return err
	}
	fmt.Print("Confirm email: ")
	_, err = fmt.Scanln(&confirmEmail)
	if err != nil {
		return err
	}

	if email != confirmEmail {
		return errors.New("Emails did not match")
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		return errors.New("A valid email address must be supplied")
	}

	password := generatePassword()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	res, err := db.Exec(`INSERT INTO users (
		name, email, password
	) VALUES (
		$1, $2, $3
	) RETURNING id`, "Admin", email, passwordHash)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return errors.New("something went wrong when creating the Admin user")
	}

	fmt.Printf("\nCreated the Admin user\nEmail: %s\nPassword: %s\n", email, password)
	fmt.Println("Please make note of these credentials as you'll need them to login.")
	fmt.Println()

	return nil
}

//GetConfig parses the lotus enviroment variables and returns the Server struct
//TODO rename me
func GetConfig() Server {
	err := godotenv.Load()
	if err != nil {
		log.Println("Couln't load .env file")
	}

	var server = Server{
		Env:      os.Getenv("LOTUS_ENV"),
		HTTPAddr: os.Getenv("LOTUS_HTTP_ADDR"),
		HTTPPort: os.Getenv("LOTUS_HTTP_PORT"),
		DBUser:   os.Getenv("LOTUS_DB_USER"),
		DBPass:   os.Getenv("LOTUS_DB_PASS"),
		DBAddr:   os.Getenv("LOTUS_DB_ADDR"),
		DBName:   os.Getenv("LOTUS_DB_NAME"),
		DB:       nil,
	}

	db, err := sql.Open("postgres", server.PostgresURL())
	if err != nil {
		log.Fatal("Failed to connect to postgres: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Couldn't ping postgres: " + err.Error())
	}

	//Delete all data in the TEST database
	if server.Env == "test" {
		if !strings.HasSuffix(server.DBName, "_test") {
			log.Fatal("Cannot run server tests on a database with a name that doesn't have the '_test' suffix. Exiting.")
		} else {
			fmt.Println("Dropping all data")
			_, err := db.Exec(fmt.Sprintf("DROP OWNED BY %s;", server.DBUser))
			if err != nil {
				log.Fatal("Error creating tables: " + err.Error())
			}
		}
	}

	//Create tables
	if server.Env == "dev" || server.Env == "test" {
		f, err := os.Open("tables.sql")
		if err != nil {
			log.Fatal("Couldn't read tables.sql: " + err.Error())
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("Error reading from file: " + err.Error())
		}

		_, err = db.Exec(string(b))
		if err != nil {
			log.Fatal("Error creating tables: " + err.Error())
		}
	}

	//Seed TEST database
	if server.Env == "test" {
		if !strings.HasSuffix(server.DBName, "_test") {
			log.Fatal("Cannot run server tests on a database with a name that doesn't have the '_test' suffix. Exiting.")
		} else {
			fmt.Println("Seeding data")
			f, err := os.Open("seed.sql")
			if err != nil {
				log.Fatal("Couldn't read seed.sql: " + err.Error())
			}

			b, err := ioutil.ReadAll(f)
			if err != nil {
				log.Fatal("Error reading from file: " + err.Error())
			}

			_, err = db.Exec(string(b))
			if err != nil {
				log.Fatal("Error seeding tables: " + err.Error())
			}
		}
	} else {
		//If there are no users, set one up
		rows := db.QueryRow(`SELECT id FROM users`)

		var userExists int
		err = rows.Scan(&userExists)
		if err == sql.ErrNoRows {
			fmt.Println("No users detected, setting up admin.")
			for {
				err := firstTimeSetup(db)
				if err != nil {
					fmt.Println(err)
				} else {
					break
				}
			}
		} else if err != nil {
			log.Fatal("Error detecting initial user: ", err.Error())
		}
	}

	server.DB = db

	return server
}
