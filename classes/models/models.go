package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)


type Class struct {
	Id 		   string
	Name       string
	Start_Date string 
	End_Date   string
	Capacity   int
}

type Classes []*Class

func (c *Classes) ToJson(w io.Writer) error{
	e:=json.NewEncoder(w)
	return e.Encode(c)
}
func DateToInt(t1,t2 time.Time) int{
	x:=t2.Sub(t1).Hours()/24
	return int(x)
}

func StringToDate(date string)time.Time{
	d,err:=time.Parse("2006-01-02 15:04:05 -0700 MST", date)
		if err!=nil{
			errors.New("error parsing to date")
		}
		return d
	}
