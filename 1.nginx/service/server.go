package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type DateResponseItem struct {
	Year  uint `json:"year"`
	Month uint `json:"month"`
	Day   uint `json:"day"`
}

type NameResponseItem struct {
	Name string `json:"name"`
}

type DateResponse []DateResponseItem
type NameResponse []NameResponseItem

const RESPONSE_ITEMS_SIZE = 10000

func date(w http.ResponseWriter, req *http.Request) {
	log.Println("json1")

	nowDate := time.Now()

	res := make([]DateResponseItem, RESPONSE_ITEMS_SIZE)
	for i := 0; i < RESPONSE_ITEMS_SIZE; i++ {
		res[i] = DateResponseItem{
			Year:  uint(nowDate.Year()),
			Month: uint(nowDate.Month()),
			Day:   uint(nowDate.Day()),
		}
	}

	resMarshalled, _ := json.Marshal(res)
	fmt.Fprintf(w, string(resMarshalled))
}

func name(w http.ResponseWriter, req *http.Request) {
	log.Println("name")

	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		name := req.PostForm.Get("name")

		res := make([]NameResponseItem, RESPONSE_ITEMS_SIZE)
		for i := 0; i < RESPONSE_ITEMS_SIZE; i++ {
			res[i] = NameResponseItem{
				Name: name,
			}
		}

		resMarshalled, _ := json.Marshal(res)
		fmt.Fprintf(w, string(resMarshalled))
	}
}

func main() {
	portToListen := os.Getenv("SERVER_PORT")

	http.HandleFunc("/date", date)
	http.HandleFunc("/name", name)

	fmt.Printf("Server is listening %s...\n", portToListen)
	http.ListenAndServe(":"+portToListen, nil)
}
