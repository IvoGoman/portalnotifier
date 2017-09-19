package login

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"github.com/IvoGoman/portalnotifier/util"

	"gopkg.in/xmlpath.v2"
)

var httpclient = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type portalClient struct {
	http        http.Client
	jar         CookieJar
	config      map[string]string
	loginTicket string
	execution   string
}

const loginticketXPath string = "//*[@id='lt']/@value"
const executionXPath string = "//*[@id='execution']/@value"
const examXPath = "//*[contains(text(), 'Prüfungen')]/@href"
const torXPath string = "//*[contains(text(), 'Noten')]/@href"
const examTableXPath = "//table[4]"
const examTableHeaderRow = "//tr[@bgcolor='#003366' and 2]"
const examTableGradeRowXPath = "//tr[@bgcolor='#EFEFEF']"

// DoLogin logs the http client into cas
func GetGrades(config map[string]string) map[string]util.Module {
	client := new(portalClient)
	client.config = config
	client.loginPortal()
	grades := client.crawlPortal()
	client.getPage(config["logout"], false)
	return grades
}

// Method retrieves necessary data for the portal login
func (client *portalClient) loginPortal() {
	config := client.config
	// retrieve login information from cas login page
	xml := client.getPage(config["server"], true)
	client.loginTicket = client.getData(xml, loginticketXPath)
	client.execution = client.getData(xml, executionXPath)
	// post authentication
	client.authenticate()
}

func (client *portalClient) crawlPortal() map[string]util.Module {
	config := client.config
	req, _ := http.NewRequest("GET", config["server"], nil)
	values := req.URL.Query()
	values.Add("service", client.config["service"])
	req.URL.RawQuery = values.Encode()
	req.Header.Add("Cookie", client.jar.Encode())
	req.Header.Add("User-Agent", client.config["user-agent"])
	res, err := client.http.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)
	portalpage := string(bytes)
	// get the url to the Exam Page and retrieve the page
	examurl := client.getData(portalpage, examXPath)
	examPage := client.getPage(examurl, false)
	// get the URL to the grade page and retrieve exam page
	torurl := client.getData(examPage, torXPath)
	torPage := client.getPage(torurl, false)
	// get the html table for grades
	table := client.getDataNode(torPage, examTableXPath)
	// extract the grades from the html table
	grades := client.processTable(table)
	return grades
}

// Method to download a page with or without cookies
func (client *portalClient) getPage(url string, cookie bool) (htmlString string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	if cookie {
		req.Header.Add("Cookie", client.jar.Encode())
	}
	req.Header.Add("User-Agent", client.config["user-agent"])

	page, err := client.http.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	client.jar.Update(page.Cookies())
	if err != nil {
		log.Fatal(err)
	}
	defer page.Body.Close()
	html, err := ioutil.ReadAll(page.Body)
	htmlString = string(html)
	return
}

// Given an HTML String and an XPath Expression the matching data is returned
func (client *portalClient) getData(htmlString string, xpath string) (val string) {
	xml, err := xmlpath.ParseHTML(strings.NewReader(htmlString))
	if err != nil {
		log.Fatal(err)
	}
	path := xmlpath.MustCompile(xpath)
	if value, ok := path.String(xml); ok {
		val = value
	}
	return
}

// Given an HTML String and an XPath Expression the matching node is returned
func (client *portalClient) getDataNode(htmlString string, xpath string) (val *xmlpath.Iter) {
	xml, err := xmlpath.ParseHTML(strings.NewReader(htmlString))
	if err != nil {
		log.Fatal(err)
	}
	path := xmlpath.MustCompile(xpath)
	val = path.Iter(xml)
	return
}

// Mirrors the actual authentication process on the Portal2 Homepage via CAS
func (client *portalClient) authenticate() {
	config := client.config
	data := url.Values{}
	data.Add("username", config["username"])
	data.Add("password", config["password"])
	data.Add("lt", client.loginTicket)
	data.Add("execution", client.execution)
	data.Add("_eventId", "submit")
	data.Add("submit", "Anmelden")
	params := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", config["server"], params)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("User-Agent", client.config["user-agent"])
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Cookie", client.jar.Encode())
	resp, err := client.http.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	client.jar.Update(resp.Cookies())
}

// Process the Table containing all the Grades
func (client *portalClient) processTable(table *xmlpath.Iter) (grades map[string]util.Module) {
	grades = make(map[string]util.Module)
	if table.Next() {
		node := table.Node()
		path := xmlpath.MustCompile(examTableGradeRowXPath)
		iter := path.Iter(node)
		grades = client.processRow(iter, grades)
	}
	return
}

// Process a single row in the grades table and extract all fields
func (client *portalClient) processRow(rows *xmlpath.Iter, grades map[string]util.Module) map[string]util.Module {
	if rows.Next() {
		rowstring := rows.Node().String()
		rowstring = strings.Replace(rowstring, "//-->", "", -1)
		rowstring = strings.Replace(rowstring, "<!--", "", -1)
		rowstring = strings.Replace(rowstring, "\t", "", -1)
		rowstring = strings.TrimSpace(rowstring)
		pattern := regexp.MustCompile("(\\s{2})+")
		rowstring = pattern.ReplaceAllString(rowstring, ", ")
		split := strings.Split(rowstring, ", ")
		var exam util.Module

		exam.ExamID, _ = strconv.ParseInt(strings.TrimSpace(split[0]), 10, 64)
		exam.Semester = split[1]
		exam.TryCountExam, _ = strconv.ParseInt(strings.TrimSpace(split[2]), 10, 64)
		exam.Date = split[3]
		exam.Name = split[4]
		exam.Prof = split[5]
		exam.Form = split[6]
		exam.Grade, _ = strconv.ParseFloat(strings.Replace(strings.TrimSpace(split[7]), ",", ".", -1), 64)
		patternBonus := regexp.MustCompile("[0-9]+.0")
		bonus := patternBonus.FindAllString(split[8], 1)[0]
		exam.Bonus, _ = strconv.ParseFloat(strings.TrimSpace(bonus), 64)
		exam.Status = split[9]
		exam.TryCountStudent = split[10]
		grades[exam.Name] = exam
		client.processRow(rows, grades)
	}
	return grades
}
