package main

import (
	"fmt"
	"os"
	"time"
)

// REST query made
// curl http://localhost:3000/GetTime/Europe/Berlin
// Time response:YYYY,MM,DD,HH,MM,SS
// 2016,04,18,15,56,08

const shortForm = "2006,01,02,15,04,05"

func timeForTZ(tz string) (string, error) {
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
	timeZone := "Europe/Berlin"
	if len(os.Args) > 1 {
		timeZone = os.Args[1]
	}

	//loc, err := time.LoadLocation("Europe/Berlin")
	timeVal, err := timeForTZ(timeZone)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Println(timeVal)
}
