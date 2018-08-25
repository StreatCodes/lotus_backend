package main

import "os"

//Server is our global struct which contains configuration information
type Server struct {
	env     string
	address string
	port    string
}

//GetConfig parses the lotus enviroment variables and returns the Server struct
func GetConfig() Server {
	var server = Server{
		env:     os.Getenv("LOTUS_ENV"),
		address: os.Getenv("LOTUS_ADRESS"),
		port:    os.Getenv("LOTUS_PORT"),
	}

	return server
}
