package main

import (
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var loginTemplate = template.Must(template.ParseFiles("templates/login.html"))
var registerTemplate = template.Must(template.ParseFiles("templates/register.html"))
var resetTemplate = template.Must(template.ParseFiles("templates/reset.html"))
var setPasswordTemplate = template.Must(template.ParseFiles("templates/set_password.html"))

type loginPostfixFn func(databaseUser)

func newLoginHandler(postfix loginPostfixFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			email := r.FormValue("email")
			password := r.FormValue("password")

			user, ok := sessionDb.getUser(email)

			if !ok || bcrypt.CompareHashAndPassword(user.password, []byte(password)) != nil {
				loginTemplate.Execute(w, map[string]interface{}{
					"Message": "Invalid email or password",
				})
				return

			} else {

				go postfix(user)

				fmt.Fprintf(w, "Welcome, %s!\n", email)
				return

			}
		}

		loginTemplate.Execute(w, nil)
	}
}

type tokenSendPostfixFn func(string)

func newRegisterHandler(postfix tokenSendPostfixFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			email := r.FormValue("email")
			if email != "" {

				_, ok := sessionDb.getUser(email)
				if ok {
					registerTemplate.Execute(w, map[string]interface{}{
						"Message": "User already exists",
					})
					return
				}

				go postfix(email)

				registerTemplate.Execute(w, map[string]interface{}{
					"Message": "Password set token sent to " + email,
				})
				return
			} else {
				registerTemplate.Execute(w, map[string]interface{}{
					"Message": "Invalid email",
				})
				return
			}
		}

		registerTemplate.Execute(w, nil)
	}
}

func newRestoreHandler(postfix tokenSendPostfixFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			email := r.FormValue("email")
			if email != "" {

				_, ok := sessionDb.getUser(email)
				if !ok {
					resetTemplate.Execute(w, map[string]interface{}{
						"Message": "User not found",
					})
					return
				}

				go postfix(email)

				resetTemplate.Execute(w, map[string]interface{}{
					"Message": "Password reset token sent to " + email,
				})
				return
			} else {
				resetTemplate.Execute(w, map[string]interface{}{
					"Message": "Invalid email",
				})
				return
			}
		}

		resetTemplate.Execute(w, nil)
	}
}

type setPasswordPostfixFn func(string, string)

func newSetPasswordHandler(postfix setPasswordPostfixFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			password := r.Form.Get("password")
			email := r.Form.Get("email")

			postfix(email, password)

			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			token := r.URL.Query().Get("token")
			u, ok := tokenServer.verifyToken(token)

			if token == "" {
				http.Error(w, "Token missing", http.StatusBadRequest)
				return
			}

			if !ok {
				http.Error(w, "Wrong token", http.StatusBadRequest)
				return
			}

			err := setPasswordTemplate.Execute(w, map[string]interface{}{
				"Username": u,
			})

			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
}
