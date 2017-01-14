package util

import (
	"io/ioutil"
	"log"
	"net/smtp"

	yaml "gopkg.in/yaml.v2"
)

func SendMail(config map[string]string, moduleMap map[string]Module) (status bool) {
	mailAuth := smtp.PlainAuth("", config["mailfrom"], config["password"], config["mailserver"])
	mailTo := []string{config["mailto"]}
	msg := []byte("From: " + config["mailfrom"] + "\r\n" +
		"To: " + config["mailto"] + "\r\n" +
		"Subject: You have new Grades\r\n" +
		"\r\n" +
		"Hello World\r\n")
	err := smtp.SendMail(config["mailserver"]+":"+config["port"], mailAuth, config["mailFrom"], mailTo, msg)
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
