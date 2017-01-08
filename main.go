package main

import (
	"fmt"
	"io"
	"strconv"

	"net/http"

	"github.com/ivogoman/portalnotifier/database"
	login "github.com/ivogoman/portalnotifier/login"
	"github.com/ivogoman/portalnotifier/util"
)

type cfg map[string]string

func grades(res http.ResponseWriter, req *http.Request) {
	gradesKnown := database.SelectGrades()
	// sort.Sort(ByName(gradesKnown))
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	response := `<DOCTYPE html>
		<html>
			<head>
				<title> Grades </title>
			</head>
			<body>
				<table>
					<tr>
						<th> Module </th>
						<th> Grade </th>
						<th> Bonus </th>
					</tr>`
	for _, grade := range gradesKnown {
		response += `<tr>
						<td>` + grade.Name + `</td>
						<td>` + strconv.FormatFloat(grade.Grade, 'f', 2, 64) + `</td>
						<td>` + strconv.FormatFloat(grade.Bonus, 'f', 0, 64) + `</td>`
	}
	response += `<tr>
				<th></th>
				<th>Average</th>
				<th>` + strconv.FormatFloat(util.CalculateAverage(gradesKnown), 'f', 2, 64) + `</th>
				</table>
				</body>
				</html>`
	io.WriteString(
		res,
		response)

}

func main() {
	database.CreateDB()
	http.HandleFunc("/grades", grades)
	http.ListenAndServe(":8080", nil)
	gradesKnown := database.SelectGrades()
	grades := login.GetGrades("./config.yml")
	for k := range gradesKnown {
		delete(grades, k)
	}
	if len(grades) > 0 {
		fmt.Println("There are " + strconv.Itoa(len(grades)) + " new grades")
		database.StoreGrades(grades)
		gradesKnown = database.SelectGrades()
		fmt.Println(util.CalculateAverage(gradesKnown))

	} else {
		fmt.Println("No new grades")
		fmt.Println(util.CalculateAverage(gradesKnown))
	}
}
