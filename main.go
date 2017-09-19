package main

import (
	"strconv"
	"time"

	"github.com/IvoGoman/portalnotifier/database"
	"github.com/IvoGoman/portalnotifier/login"
	"github.com/IvoGoman/portalnotifier/util"
	"github.com/IvoGoman/portalnotifier/html"
	"os"
)

var path, _ = os.Getwd()

var configPath = path + "/portalnotifier/config.yml"


var config = util.LoadConfig(configPath)


var gradesKnown = make(map[string]util.Module)

func main() {
	database.CreateDB()
	gradesKnown = database.SelectGrades()
	interval, _ := strconv.Atoi(config["interval"])
	go html.Serve(config)
	grades := login.GetGrades(config)
	checkGrades(gradesKnown, grades)
	gradeTicker := time.NewTicker(time.Minute * time.Duration(interval))
	for range gradeTicker.C {
		grades := login.GetGrades(config)
		checkGrades(gradesKnown, grades)
	}
}

func checkGrades(known map[string]util.Module, current map[string]util.Module) {
	for k := range known {
		delete(current, k)
	}
	if len(current) > 0 {
		database.StoreGrades(current)
		known = database.SelectGrades()
		util.SendMail(config, current, util.CalculateAverage(known))

	}

}
