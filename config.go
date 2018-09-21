package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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
	DB       *sqlx.DB
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

	db, err := sqlx.Open("postgres", server.PostgresURL())
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
	}

	server.DB = db

	return server
}
