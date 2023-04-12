package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	// Set up authentication information.
	smtpServer := "smtp.rambler.ru"
	smtpPort := 465
	smtpUsername := "rlebovski"
	smtpPassword := "glorytotheatP2"

	// Set up the message.
	message := []byte("To: googlodict@gmail.com\r\n" +
		"Subject: Test email\r\n" +
		"\r\n" +
		"Hello,\r\n" +
		"This is a test email sent using Go!\r\n")

	// Set up the SMTP client.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	err := smtp.SendMail(smtpAddr, auth, smtpUsername, []string{"googlodict@gmail.com"}, message)
	if err != nil {
		fmt.Printf("Error sending email: %v", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
