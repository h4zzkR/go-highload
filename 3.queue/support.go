package main

import (
	"log"
	"encoding/json"

	email "queue/email"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func MailItem2JSON(item email.MailItem) ([]byte, error) {
	request, err := json.Marshal(item)
	return request, err
}