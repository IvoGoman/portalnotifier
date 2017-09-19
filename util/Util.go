package util

import (
	"io/ioutil"
	"log"
	"net/smtp"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// Sends a mail to the user with the new grades and their average grade
func SendMail(config map[string]string, moduleMap map[string]Module, average float64) (status bool) {
	mailAuth := smtp.PlainAuth("", config["mailfrom"], config["password"], config["mailserver"])
	mailTo := []string{config["mailto"]}
	msg := "From: " + config["mailfrom"] + "\r\n" +
		"To: " + config["mailto"] + "\r\n" +
		"Subject: You have new Grades\r\n" +
		"\r\n"
	for _, grade := range moduleMap {
		msg += "You got a " + strconv.FormatFloat(grade.Grade, 'f', 2, 64) + " in " + grade.Name + "\r\n"
	}

	msg += strconv.FormatFloat(average, 'f', 2, 64) + " is your new average. \r\n"
	err := smtp.SendMail(config["mailserver"]+":"+config["port"], mailAuth, config["mailFrom"], mailTo, []byte(msg))
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

// LoadConfig loads config from file
func LoadConfig(config string) map[string]string {
	content := make(map[string]string)
	raw, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(raw, &content)
	return content
}
