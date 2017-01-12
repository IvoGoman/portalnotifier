package util

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"sort"
	"strconv"

	"github.com/IvoGoman/portalnotifier/login"

	yaml "gopkg.in/yaml.v2"
)

func SendMail(config map[string]string, moduleMap map[string]Module) (status bool) {
	mailAuth := smtp.PlainAuth("", config["mailfrom"], config["password"], config["mailserver"])
	mailTo := []string{config["mailto"]}
	msg := []byte("To: " + config["mailto"] + "\r\n" +
		"Subject: You have new Grades\r\n" +
		"\r\n" +
		"Hello World\r\n")
	err := smtp.SendMail(config["mailserver"]+":"+config["mailport"], mailAuth, config["mailFrom"], mailTo, msg)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// LoadConfig loads config from file
func LoadConfig(config string) map[string]string {
	content := make(map[string]string)
	raw, _ := ioutil.ReadFile(config)
	yaml.Unmarshal(raw, &content)
	return content
}

// Serves the Grade Table
func Serve() {
	http.HandleFunc("/grades", grades)
	http.ListenAndServe(":8080", nil)
}
func grades(res http.ResponseWriter, req *http.Request) {
	gradesCurrent := login.GetGrades(config)
	grades := ModuleMapToArray(gradesCurrent)
	sort.Sort(ByName(grades))
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
				<th>` + strconv.FormatFloat(CalculateAverage(gradesKnown), 'f', 2, 64) + `</th>
				</table>
				</body>
				</html>`
	io.WriteString(
		res,
		response)

}
