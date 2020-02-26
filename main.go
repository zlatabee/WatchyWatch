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

type Rule struct {
	Type       string `json:"type"`
	Regex      string `json:"regex"`
	IgnoreCase bool   `json:"ignore_case,omitempty"`
}

type Category []interface{}

func getTimeperiods() string {
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Error loading New York location: %v", err)
	}
	now := time.Now().In(newYork)
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, newYork)
	beginDate := endDate.Add(-24 * time.Hour)
	return beginDate.Format(time.RFC3339) + "/" + endDate.Format(time.RFC3339)
}

func getCategories() string {
	categories := []Category{
		{[]string{"Work"}, Rule{Type: "regex", Regex: "Google Docs"}},
		{[]string{"Work", "Programming"}, Rule{Type: "regex", Regex: "GitHub|Stack Overflow"}},
		{[]string{"Work", "Programming", "ActivityWatch"}, Rule{Type: "regex", Regex: "ActivityWatch|aw-", IgnoreCase: true}},
		{[]string{"Media", "Games"}, Rule{Type: "regex", Regex: "Minecraft|RimWorld"}},
		{[]string{"Media", "Video"}, Rule{Type: "regex", Regex: "YouTube|Plex"}},
		{[]string{"Media", "Social Media"}, Rule{Type: "regex", Regex: "reddit|Facebook|Twitter|Instagram", IgnoreCase: true}},
		{[]string{"Comms", "IM"}, Rule{Type: "regex", Regex: "Messenger|Telegram|Signal|WhatsApp"}},
		{[]string{"Comms", "Email"}, Rule{Type: "regex", Regex: "Gmail"}}}
	catBytes, err := json.Marshal(categories)
	if err != nil {
		log.Fatalf("Error eeeeee: %v", err)
	}
	return string(catBytes)
}

func main() {
	queryStr := `
  window_events = query_bucket(find_bucket('aw-watcher-window_'));
  not_afk_events = query_bucket(find_bucket('aw-watcher-afk_'));
  not_afk_events = filter_keyvals(not_afk_events, "status", ["not-afk"]);
  events = filter_period_intersect(window_events, not_afk_events);

  classes = ` + getCategories() + `;
  events = categorize(events, classes);
  events = merge_events_by_keys(events,["$category"]);

  RETURN = events;
  `

	query := strings.SplitAfter(queryStr, ";")
	req := Request{Query: query, Timeperiods: []string{getTimeperiods()}}
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

	if len(events) != 1 {
		log.Fatalf("Expected 'events' to be length 1, instead got: %v", events)
	}

	catMap := make(map[string]float64)

	for _, event := range events[0] {
		fillMap(catMap, event.Data.Category, event.Duration/3600.0)
	}

	for cat, dur := range catMap {
		log.Printf("%s: %0.1f hours", cat, dur)
	}
}

func fillMap(catMap map[string]float64, cat []string, dur float64) {
	catStr := strings.Join(cat, "->")
	catMap[catStr] += dur
	if len(cat) > 1 {
		fillMap(catMap, cat[:len(cat)-1], dur)
	}
}
