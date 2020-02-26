package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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
	Category  []string `json:"$category"`
}

type Event struct {
	Id        string
	Timestamp string
	Duration  float64
	Data      Data
}

type Request struct {
	Query       []string `json:"query"`
	Timeperiods []string `json:"timeperiods"`
}

func main() {
	queryStr := `
  window_events = query_bucket(find_bucket('aw-watcher-window_'));
  not_afk_events = query_bucket(find_bucket('aw-watcher-afk_'));
  not_afk_events = filter_keyvals(not_afk_events, "status", ["not-afk"]);
  events = filter_period_intersect(window_events, not_afk_events);

  classes = [[["Work"],{"type":"regex","regex":"Google Docs"}],[["Work","Programming"],{"type":"regex","regex":"GitHub|Stack Overflow"}],[["Work","Programming","ActivityWatch"],{"type":"regex","regex":"ActivityWatch|aw-","ignore_case":true}],[["Media","Games"],{"type":"regex","regex":"Minecraft|RimWorld"}],[["Media","Video"],{"type":"regex","regex":"YouTube|Plex"}],[["Media","Social Media"],{"type":"regex","regex":"reddit|Facebook|Twitter|Instagram","ignore_case":true}],[["Comms","IM"],{"type":"regex","regex":"Messenger|Telegram|Signal|WhatsApp"}],[["Comms","Email"],{"type":"regex","regex":"Gmail"}]];

  events = categorize(events, classes);
  events = merge_events_by_keys(events,["$category"]);

  RETURN = events;
  `

	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Error loading New York location: %v", err)
	}
	timeperiods := time.Now().Add(-24*time.Hour).In(ny).Format(time.RFC3339) + "/" + time.Now().In(ny).Format(time.RFC3339)

	query := strings.SplitAfter(queryStr, ";")
	req := Request{Query: query, Timeperiods: []string{timeperiods}}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Error aaaaaaaa: %v", err)
	}

	resp, err := http.Post("http://localhost:5600/api/0/query", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatalf("Error oh no: %v", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Another error: %v", err)
	}

	events := make([][]Event, 0)

	if err := json.Unmarshal(respBody, &events); err != nil {
		log.Fatalf("Uh oh: %v", err)
	}

	log.Print(events)
}
