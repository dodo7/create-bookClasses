package main

import (
	"encoding/json"
	"log"
	"net/http"
	"story/classes/database"
	"story/classes/models"

	uuid "github.com/satori/go.uuid"
)




func main() {
	database.Connect()
	mux := http.NewServeMux()
	mux.HandleFunc("/createClass", CreateClass)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func CreateClass(w http.ResponseWriter,r *http.Request)  {
	if r.Method != "POST"{
		http.Error(w,http.StatusText(405),http.StatusMethodNotAllowed)
		return
	}
	class:=models.Class{}
	err:=json.NewDecoder(r.Body).Decode(&class)
	if err!=nil{
		panic(err)
	}
		sd:=models.StringToDate(class.Start_Date)
		ed:=models.StringToDate(class.End_Date)
		x:=models.DateToInt(sd,ed)
		if class.Capacity<0 || x<=0{
			http.Error(w,"bad capacity or dates",http.StatusNotAcceptable)
			return
		}
	inputJ,err:=json.Marshal(class)
	if err!=nil{
		panic(err)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	_,err=w.Write(inputJ)
	if err!=nil{
		panic(err)
	}
	
	for i:=0;i<=x;i++{
		startDate:=models.StringToDate(class.Start_Date)
		std:=startDate.AddDate(0,0,i)
		end:=startDate.AddDate(0,0,i+1)
		cID:= uuid.NewV4().String()
		cl:=&models.Class{
		Id: cID,
		Name:       class.Name,
		Start_Date:  std.String(),
		End_Date:   end.String(),
		Capacity:   class.Capacity,
	}
	
	tx,txErr:=database.Db.Begin()
	if txErr!=nil{
		log.Printf("Error while opening tx,Reason %v\n",txErr)
	}
	_,Inserterr:=tx.Model(cl).Returning("*").Insert()
	if Inserterr!=nil{
		log.Printf("Error while inserting tx,Reason %v\n",Inserterr)
		tx.Rollback()
	}
	tx.Commit()
	log.Printf("Successfully submited")
}
	w.WriteHeader(http.StatusOK)
	
}