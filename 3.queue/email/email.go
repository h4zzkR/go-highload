package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

const (
	serverName = "smtp.gmail.com:465"
)

type MailItem struct {
	To    string
	Theme string
	Body  string
}

type MailConfig struct {
	From string
	auth smtp.Auth
	host string
}

func Setup(from string, pass string) *MailConfig {
	host, _, _ := net.SplitHostPort(serverName)
	auth := smtp.PlainAuth("", from, pass, host)
	return &MailConfig{From: from, auth: auth, host: host}
}

func Dial(to string, message string, conf *MailConfig) {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         conf.host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", serverName, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, conf.host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(conf.auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(conf.From); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
}

func Send(to []string, subj string, body string, conf *MailConfig) {
	for _, recep := range to {
		headers := make(map[string]string)
		headers["From"] = conf.From
		headers["To"] = recep
		headers["Subject"] = subj

		// Setup message
		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + body

		Dial(recep, message, conf)
	}

}

// func main() {

// 	conf := Setup("megaterraboy@gmail.com", "ixvvgvwikkguzemf")
// 	Send([]string{"googlodict@gmail.com"}, "FUCK!", "YOU.", conf)
// }

// ixvvgvwikkguzemf
// eeerelvuwkznopru
