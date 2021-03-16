package main

import (
	"net/http"
	"story/booking/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

func GetUser(w http.ResponseWriter, req *http.Request) models.User {
	// get cookie
	c, err := req.Cookie("session")
	if err != nil {
		sID:= uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	// if the user exists already, get user
	var u models.User
	if s, ok := dbSessions[c.Value]; ok {
		s.LastActivity = time.Now()
		dbSessions[c.Value] = s
		u = dbUsers[s.Un]
	}
	return u
}

func AlreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}
	s, ok := dbSessions[c.Value]
	if ok {
		s.LastActivity = time.Now()
		dbSessions[c.Value] = s
	}
	_, ok = dbUsers[s.Un]
	// refresh session
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	return ok
}

func CleanSessions() {

	for k, v := range dbSessions {
		if time.Now().Sub(v.LastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
}


