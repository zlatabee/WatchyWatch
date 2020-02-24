package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	App       string
	Url       string
	Title     string
	Type      string
	Duration  float64
	Timestamp string
	Audible   bool
	Incognito bool
}

type Event struct {
	Timestamp string
	Duration  float64
	Data      Data
}

type Bucket struct {
	Id       string
	Created  string
	Name     string
	Type     string
	Client   string
	Hostname string
	Events   []Event
}

type Export struct {
	Buckets map[string]Bucket
}

func main() {
	resp, err := http.Get("http://localhost:5600/api/0/export")
	if err != nil {
		log.Fatalf("Error oh no: %v", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Another error: %v", err)
	}

	var export Export
	if err := json.Unmarshal(respBody, &export); err != nil {
		log.Fatalf("Uh oh: %v", err)
	}
	log.Println(export.Buckets)
}
