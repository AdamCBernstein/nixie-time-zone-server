package main

// Go implementation of time-zone-server project.
// Configure the "time server" for an Ian Sparkes Nixie tube clock to
// a time zone entry, similar to the following:
//   http://mpgedit.org:3000/GetTime/Europe/Zurich
//
// This will return the time for Zurich Switzerland.
// Specify your local time zone to keep your clock time automatically set.
//
// https://github.com/isparkes/time-zone-server#time-zone-server
// This server responds to a REST call in this format:
//   http://host-location:port/GetTime/Country/City
//
// Actual example:
//   http://mpgedit.org:3000/GetTime/America/Los_Angeles
//
// Response to this call:
// 2021,12,08,18,32,18
// YEAR MM DD HH MM SS
//
// These are the supported REST calls:
//	GetTime
//	GetTimeRaw
//	GetTimeOffset
//	GetTimeZone
//
// Build this project:
//    go build
//
// Run this time server:
//    ./nixie_timeserver &
//
// The default port this server is 3000.
// This can be changed by doing the following:
//    export PORT=4444
//   ./nixie_timeserver &

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	restGetTime       = "/GetTime/"
	restGetTimeRaw    = "/GetTimeRaw/"
	restGetTimeOffset = "/GetTimeOffset/"
	restGetTimeZone   = "/GetTimeZone/"
	shortTimeFormat   = "2006,01,02,15,04,05"
	maxClients        = 10
)

var (
	httpMax  = make(chan struct{}, maxClients)
	httpPort = ":3000"
)

func getTime(w http.ResponseWriter, r *http.Request) {
	httpMax <- struct{}{}
	defer func() { <-httpMax }()

	path := r.URL.Path
	if !strings.HasPrefix(path, restGetTime) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad URL format: %s\n", path)
		return
	}

	locationTz, err := time.LoadLocation(path[len(restGetTime):])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invlid timezone specified: %s\n", path)
		return
	}

	fmt.Fprintln(w, time.Now().In(locationTz).Format(shortTimeFormat))
}

func getTimeRaw(w http.ResponseWriter, r *http.Request) {
	httpMax <- struct{}{}
	defer func() { <-httpMax }()

	path := r.URL.Path
	if !strings.HasPrefix(path, restGetTimeRaw) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad URL format: %s\n", path)
		return
	}

	locationTz, err := time.LoadLocation(path[len(restGetTimeRaw):])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invlid timezone specified: %s\n", path)
		return
	}
	fmt.Fprintln(w, time.Now().In(locationTz).Format(time.RFC1123Z))
}

func getTimeOffset(w http.ResponseWriter, r *http.Request) {
	httpMax <- struct{}{}
	defer func() { <-httpMax }()

	path := r.URL.Path
	if !strings.HasPrefix(path, restGetTimeOffset) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad URL format: %s\n", path)
		return
	}

	locationTz, err := time.LoadLocation(path[len(restGetTimeOffset):])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invlid timezone specified: %s\n", path)
		return
	}

	fmt.Fprintln(w, "GMT"+time.Now().In(locationTz).Format("-0700"))
}

func getTimeZone(w http.ResponseWriter, r *http.Request) {
	httpMax <- struct{}{}
	defer func() { <-httpMax }()

	path := r.URL.Path
	if !strings.HasPrefix(path, restGetTimeZone) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad URL format: %s\n", path[len(restGetTimeZone):])
		return
	}

	locationTz, err := time.LoadLocation(path[len(restGetTimeZone):])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invlid timezone specified: %s\n", path)
		return
	}

	zoneForUrl, _ := time.Now().In(locationTz).Zone()
	fmt.Fprintln(w, zoneForUrl)
}

func main() {
	port := os.Getenv("PORT")
	if len(port) > 0 {
		httpPort = ":" + port
	}
	http.HandleFunc(restGetTime, getTime)
	http.HandleFunc(restGetTimeRaw, getTimeRaw)
	http.HandleFunc(restGetTimeOffset, getTimeOffset)
	http.HandleFunc(restGetTimeZone, getTimeZone)
	log.Fatal(http.ListenAndServe(httpPort, nil))
}
