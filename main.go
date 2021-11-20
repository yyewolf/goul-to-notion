package main

import (
	"time"

	"github.com/yyewolf/notionapi"
)

var client *notionapi.Client

func main() {
	client = notionapi.NewClient(notionapi.Token(notionSecret))

	dewIt()

	year, month, day := time.Now().Date()
	hour := 20
	delay := time.Until(time.Date(year, month, day, hour, 0, 0, 0, time.Local))
	if delay < 0 {
		delay = time.Until(time.Date(year, month, day+1, hour, 0, 0, 0, time.Local))
	}
	time.Sleep(delay)
	t := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-t.C:
				dewIt()
			}
		}
	}()
}

func dewIt() {
	emptyLines()
	// We check what time it is
	lastMonday := time.Now()
	nextSunday := time.Now()
	for lastMonday.Weekday() != time.Monday {
		lastMonday = lastMonday.AddDate(0, 0, -1)
	}
	for nextSunday.Weekday() != time.Sunday {
		nextSunday = nextSunday.AddDate(0, 0, 1)
	}

	addLines(lastMonday, nextSunday)
}
