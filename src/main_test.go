package main

import (
	"fmt"
)

var server Server

func init() {
	//Setup DB
	fmt.Println("Running Setup")
	server = GetConfig()
	fmt.Println()
}
