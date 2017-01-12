package main

import (
	"fmt"
	"strconv"

	"github.com/IvoGoman/portalnotifier/html"
	login "github.com/IvoGoman/portalnotifier/login"
	"github.com/IvoGoman/portalnotifier/util"
)

type cfg map[string]string

var config = util.LoadConfig("./config.yml")

var gradesKnown = login.GetGrades(config)

// var gradesKnown = make(map[string]util.Module)

func main() {
	// database.CreateDB()

	html.Serve(config)
	grades := login.GetGrades(config)
	for k := range gradesKnown {
		delete(grades, k)
	}
	if len(grades) > 0 {
		fmt.Println("There are " + strconv.Itoa(len(grades)) + " new grades")
		// database.StoreGrades(grades)
		// status := util.SendMail(config, grades)

		// if status == false {
		// 	fmt.Println("Mail not send")
		// } else {
		// 	fmt.Println("Mail send")
		// }

		gradesKnown = login.GetGrades(config)
		fmt.Println(util.CalculateAverage(gradesKnown))

	} else {
		fmt.Println("No new grades")
		fmt.Println(util.CalculateAverage(gradesKnown))
	}
}
