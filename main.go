package main

import (
	"strconv"
	"time"

	"fmt"

	"github.com/IvoGoman/portalnotifier/database"
	"github.com/IvoGoman/portalnotifier/login"
	"github.com/IvoGoman/portalnotifier/util"
	"github.com/ivogoman/portalnotifier/html"
)

type cfg map[string]string

var config = util.LoadConfig("./config.yml")

var gradesKnown = make(map[string]util.Module)

// var gradesKnown = make(map[string]util.Module)

func main() {
	database.CreateDB()
	interval, _ := strconv.Atoi(config["interval"])
	go html.Serve(config)
	grades := login.GetGrades(config)
	gradeTicker := time.NewTicker(time.Minute * time.Duration(interval))
	for t := range gradeTicker.C {
		fmt.Println(t)
		for k := range gradesKnown {
			delete(grades, k)
		}
		if len(grades) > 0 {
			database.StoreGrades(grades)
			gradesKnown = login.GetGrades(config)
			util.SendMail(config, grades)
		}
	}
}
