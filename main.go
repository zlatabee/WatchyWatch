package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:5600/api/0/info")
	if err != nil {
		log.Fatalf("Error oh no: %v", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Another error: %v", err)
	}
	log.Printf("Server info: %s", respBody)
	log.Println("Hi")
}
