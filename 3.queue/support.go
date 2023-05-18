package main

import (
	"encoding/json"
	"fmt"
	"log"

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

func helpInfo() {
	msg := `Welcome! Here some useful information:
- Visit /register for creating your brand new account. All you need is to input email,
  and receive message with token link. By following the link you can set password for new account.
- Visit /login for, emm, login.
- Forgot password, silly you? We'll forgive you this time, just visit /login and click "Restore password" link.
  Follow the link in email message and set your new password. Try to remember it, though.
- Pay attention - tokens have their lifetimes, so hurry on to do your things.
----------------------------------------------------------------------------------------------------------------------------------------------
- LIMITATIONS:
  Although the invalidation of UNUSED tokens works after a while, there is no check in the system for the invalidation of already used tokens.
  It's because I don't have time to do it. Well, because there was no technical task.`
	fmt.Print(msg, "\n")
}
