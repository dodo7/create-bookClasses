package models

import (
	"time"
)


type User struct {
	UserName string
	Password []byte
	First    string
	Last     string
	Role     string
}

type Session struct {
	Un           string
	LastActivity time.Time
}
type Booking struct{
	Username string
	StartingDate string
	EndingDate string
	ClassID string
}
type Class struct {
	Name       string
	Start_Date string 
	End_Date   string
	Capacity   int
	Id		string `pg:",pk"`
}



