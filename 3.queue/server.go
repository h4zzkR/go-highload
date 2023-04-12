package main

import (
	"fmt"
	"net/http"
)

/**
DISCLAIMER - GOVNOCODE
- bare http
- map instead normal database
*/

const (
	addr      = "127.0.0.1:8080"
	mQAddr    = "amqp://guest:guest@localhost:5672/"
	hostname  = "ApertureSciencePortal"
	jwtSecret = "secretTODO"
)

var sessionDb = newDatabase()
var emailQueue = newMessageQueue()
var tokenServer = newTokenServer(jwtSecret)

func main() {

	mux := http.NewServeMux()

	defer emailQueue.conn.Close()
	defer emailQueue.ch.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Proudly served with Go and HTTPS!\n")
	})

	mux.HandleFunc("/login", newLoginHandler(greetLoginPostfix))
	mux.HandleFunc("/register", newRegisterHandler(registerPostfix))
	mux.HandleFunc("/restore", newRestoreHandler(registerPostfix))
	mux.HandleFunc("/set-password", newSetPasswordHandler(setPasswordPostfix))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("Starting server on port 8080...")
	err := srv.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
