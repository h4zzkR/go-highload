package main

import (
	"fmt"
	email "queue/email"
)

func greetLoginPostfix(user databaseUser) {
	mail := email.MailItem{
		To:    user.email,
		Theme: "New login detected",
		Body:  "Achtung!",
	}

	emailQueue.sendEmail(mail)
}

func registerPostfix(userEmail string) {
	tkn, _ := tokenServer.generateToken(userEmail)

	link := fmt.Sprintf("http://%s/set-password?token=%s", addr, tkn)

	mail := email.MailItem{
		To:    userEmail,
		Theme: fmt.Sprintf("Welcome to %s!", hostname),
		Body:  fmt.Sprintf("Click on this link to set your new pasword: %s", link),
	}

	emailQueue.sendEmail(mail)
}

func setPasswordPostfix(email, password string) {
	sessionDb.adddatabaseUser(email, password)
}
