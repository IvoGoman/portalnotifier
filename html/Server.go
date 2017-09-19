package html

import (
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/IvoGoman/portalnotifier/database"
	"github.com/IvoGoman/portalnotifier/util"
)

var config = make(map[string]string)

// Serves the Grade Table
func Serve(cfg map[string]string) {
	config = cfg
	http.HandleFunc("/grades", grades)
	http.ListenAndServe(":"+config["serverport"], nil)
}

// Creates the plain html table that represents the grades
func grades(res http.ResponseWriter, req *http.Request) {
	gradesKnown := database.SelectGrades()
	grades := util.ModuleMapToArray(gradesKnown)
	sort.Sort(util.ByName(grades))
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
	for _, grade := range grades {
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
