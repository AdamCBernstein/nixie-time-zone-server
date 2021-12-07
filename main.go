package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	restGetTime     = "/getTime/"
	shortTimeFormat = "2006,01,02,15,04,05"
	maxClients      = 10
	httpPort        = ":8082"
)

func startRestTzServer() {
	httpMax := make(chan struct{}, maxClients)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpMax <- struct{}{}
		defer func() { <-httpMax }()

		if httpGetUrl := html.EscapeString(r.URL.Path); len(httpGetUrl) == 0 {
			fmt.Fprintf(w, "Bad URL format: %s\n", httpGetUrl)
		} else if offset := strings.Index(httpGetUrl, restGetTime); offset < 0 {
			fmt.Fprintf(w, "Bad REST URI: %s", httpGetUrl)
		} else if timeForURL, err := timeForTz(httpGetUrl[offset+len(restGetTime):]); err != nil {
			fmt.Fprintf(w, "%q No such time zone %v", httpGetUrl[offset+len(restGetTime):], err)
		} else {
			fmt.Fprintf(w, "%s", timeForURL)
		}
	})
	log.Fatal(http.ListenAndServe(httpPort, nil))
}

func timeForTz(tz string) (string, error) {
	locationTz, err := time.LoadLocation(tz)
	if err != nil {
		return "", err
	}

	n := time.Now().In(locationTz)
	nowTime := fmt.Sprintf("%4d,%02d,%02d,%02d,%02d,%02d",
		n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())
	if tzTime, err := time.ParseInLocation(shortTimeFormat, nowTime, locationTz); err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		return "", err
	} else {
		timeFmt := tzTime.Format(shortTimeFormat)
		return timeFmt, nil
	}
}

func main() {
	startRestTzServer()
}
