package main

import (
	"strconv"
	"time"

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
	gradesKnown = database.SelectGrades()
	interval, _ := strconv.Atoi(config["interval"])
	go html.Serve(config)
	gradeTicker := time.NewTicker(time.Millisecond * time.Duration(interval))
	for t := range gradeTicker.C {
		grades := login.GetGrades(config)
		for k := range gradesKnown {
			delete(grades, k)
		}
		if len(grades) > 0 {
			database.StoreGrades(grades)
			gradesKnown = database.SelectGrades()
			util.SendMail(config, grades, util.CalculateAverage(gradesKnown))
		}
	}
}
