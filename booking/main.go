package main

import (
	"encoding/json"
	"log"
	"net/http"
	"story/booking/database"
	"story/booking/models"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)
var dbUsers = map[string]models.User{}       // user ID, user
var dbSessions = map[string]models.Session{}
var dbSessionsCleaned time.Time
const sessionLength int = 30

func main() {
	database.Connect()
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/bookclass", BookClass)
	mux.HandleFunc("/signup", Signup)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Authorized(Logout))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8080", mux))

	
}

func Index(w http.ResponseWriter, req *http.Request) {
	 GetUser(w, req)
	
}

func BookClass(w http.ResponseWriter,r *http.Request)  {
	u := GetUser(w, r)
	if !AlreadyLoggedIn(w, r) {
		http.Error(w, "Not Logged in", http.StatusBadRequest)
		return
	}
	//check if the person has role user, if not cant book
	if u.Role != "user" {
		http.Error(w, "You must be user to enter the book class", http.StatusForbidden)
		return
	}else if u.Role=="user"{
		b:=models.Booking{}
		c:=models.Class{}
		err:=json.NewDecoder(r.Body).Decode(&b)
	if err!=nil{
		panic(err)
	}
	//receive and check if the class has enough capacity
	var ca int
	database.Db.Model(&c).Column("class.capacity").Where("id=?",b.ClassID).Select(&ca)
	
	
		if ca==0  {
		http.Error(w, "You must enter an existing class name or class is full", http.StatusBadRequest)
		return
	}

	//receive the start and end date fro making the booking for the class had been choosen
	var sd,ed string
	err=database.Db.Model(&c).Column("class.start__date","class.end__date").Where("id=?",b.ClassID).Select(&sd,&ed)
	if err!=nil{
		panic(err)
	}
		//insert on bookings
	bc:=models.Booking{
		Username:     u.UserName,
		StartingDate: sd,
		EndingDate: ed,
		ClassID:      b.ClassID,
	}
		//insert on db the booking
	_,err=database.Db.Model(&bc).Returning("*").Insert()
	if err!=nil{
		panic(err)
	}
	//set capacity on db after we have succesfully make a new booking
	_,err=database.Db.Model(&c).Set("capacity=?",ca-1).Where("id=?",b.ClassID).Update()
	if err!=nil{
		panic(err)
	}
	
	w.WriteHeader(http.StatusOK)
	}

}

func Signup(w http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u models.User
	// process form submission
	if req.Method == http.MethodPost {
		err:=json.NewDecoder(req.Body).Decode(&u)
	if err!=nil{
		panic(err)
	}
	us:=models.User{
		UserName: u.UserName,
		Password: u.Password,
		First:    u.First,
		Last:    u.Last,
		Role:     u.Role,
	}
		// username taken?
		if _, ok := dbUsers[us.UserName]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// create session
		sID:= uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = models.Session{us.UserName, time.Now()}
		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		u = models.User{us.UserName, bs, us.First, us.Last, us.Role}
		dbUsers[us.UserName] = u
		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	
}

func Login(w http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u models.User
	// process form submission
	if req.Method == http.MethodPost {
		err:=json.NewDecoder(req.Body).Decode(&u)
	if err!=nil{
		panic(err)
	}
	us:=models.User{
		UserName: u.UserName,
		Password: u.Password,
	}
		// is there a username?
		u, ok := dbUsers[us.UserName]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// does the entered password match the stored password?
		err = bcrypt.CompareHashAndPassword(u.Password, []byte(us.Password))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		sID:= uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = models.Session{u.UserName, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	
}

func Logout(w http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("session")
	// delete the session
	delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go CleanSessions()
	}

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func Authorized(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// code before
		if !AlreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
		// code after
	})
}



