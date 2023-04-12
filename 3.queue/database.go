package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type databaseUser struct {
	password []byte
	email    string
}

type database struct {
	storage map[string]databaseUser
}

func (d *database) adddatabaseUser(email, password string) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("bcrypt error")
	}
	d.storage[email] = databaseUser{
		password: hashedPass,
		email:    email,
	}
}

func (d *database) getUser(uname string) (databaseUser, bool) {
	u, ok := d.storage[uname]
	return u, ok
}

func (d *database) getHashedPass(uname string) ([]byte, bool) {
	u, ok := d.storage[uname]
	return u.password, ok
}

func newDatabase() *database {
	db := database{
		storage: make(map[string]databaseUser),
	}
	db.adddatabaseUser("bedrock@bedrock.com", "bedrock")
	return &db
}
