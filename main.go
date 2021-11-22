package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yyewolf/notionapi"
)

var client *notionapi.Client

func main() {
	client = notionapi.NewClient(notionapi.Token(notionSecret))

	dewIt()
	doLoop()

	time1, time2 := 18, 7

	year, month, day := time.Now().Date()
	hour := time1
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

	time.Sleep(time.Duration(24-time1+time2) * time.Hour)
	dewIt()
	t2 := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-t2.C:
				dewIt()
			}
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
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
