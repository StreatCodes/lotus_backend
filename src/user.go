package main

import "time"

//User contains all the information required to manipulate
//the administrator users in Lotus
type User struct {
	ID        int
	CreatedAt time.Time
	Name      string
	Email     string
}
