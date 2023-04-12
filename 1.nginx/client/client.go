package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const SERVER = "http://localhost:8080"

/*
*
Client code for POST-requests.
*/
func main() {
	data := url.Values{
		"name": {"John Doe"},
	}

	resp, err := http.PostForm(SERVER+"/name", data)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	var res []map[string]string
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res)
}
