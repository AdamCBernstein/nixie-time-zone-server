package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

const (
	shortTimeFormat = "2006,01,02,15,04,05"
	maxClients      = 10
	httpPort        = ":8081"
)

func startRestTzServer() {
	httpMax := make(chan struct{}, maxClients)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpMax <- struct{}{}
		defer func() { <-httpMax }()

		if timeStr, err := timeForTz(html.EscapeString(r.URL.Path)[1:]); err != nil {
			fmt.Fprintf(w, "No such time zone %q", html.EscapeString(r.URL.Path))
		} else {
			fmt.Fprintf(w, "%s", timeStr)
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
