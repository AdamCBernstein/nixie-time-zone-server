package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

const shortForm = "2006,01,02,15,04,05"

func startRestTzServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if timeStr, err := timeForTz(html.EscapeString(r.URL.Path)[1:]); err != nil {
			fmt.Fprintf(w, "No such time zone %q", html.EscapeString(r.URL.Path))
		} else {
			fmt.Fprintf(w, "%s", timeStr)
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func timeForTz(tz string) (string, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return "", err
	}

	n := time.Now().In(loc)
	nstr := fmt.Sprintf("%4d,%02d,%02d,%02d,%02d,%02d",
		n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())
	if t, err := time.ParseInLocation(shortForm, nstr, loc); err != nil {
		fmt.Println(fmt.Errorf("%v", err))
		return "", err
	} else {
		timeFmt := t.Format(shortForm)
		return timeFmt, nil
	}
}

func main() {
	startRestTzServer()
}
