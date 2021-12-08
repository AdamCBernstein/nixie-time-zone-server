package main

import (
	"fmt"
	"log"
	"net/http"
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
	httpPort          = ":3000"
)

var (
	httpMax = make(chan struct{}, maxClients)
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
	http.HandleFunc(restGetTime, getTime)
	http.HandleFunc(restGetTimeRaw, getTimeRaw)
	http.HandleFunc(restGetTimeOffset, getTimeOffset)
	http.HandleFunc(restGetTimeZone, getTimeZone)
	log.Fatal(http.ListenAndServe(httpPort, nil))
}
