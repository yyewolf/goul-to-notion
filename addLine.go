package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yyewolf/notionapi"
)

var column = []string{"Heure", "Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi"}

func timeToMin(t time.Time) int {
	return 60*t.Hour() + t.Minute()
}

func emptyLines() {
	resp, _ := client.Search.Do(context.Background(), &notionapi.SearchRequest{
		PageSize: 30,
	})
	for _, r := range resp.Results {
		if fmt.Sprint(r.GetObject()) == "page" {
			d, err := json.Marshal(r)
			if err != nil {
				continue
			}
			p := &notionapi.Page{}
			json.Unmarshal(d, p)
			client.Block.Delete(context.Background(), notionapi.BlockID(p.ID))
		}
	}
}

func addLines(from, to time.Time) {
	table := getTablePrimitive(token, from.Unix(), to.Unix())
	if len(table.Data.Table.Plannings) <= 0 {
		fmt.Println("table empty")
		return
	}
	events := table.Data.Table.Plannings[0].Events
	firstDay := events[0].Start.Day()
	lastDay := events[len(events)-1].End.Day()
	var lines [][]tableEvent
	var space = 30 * time.Minute
	var earliest = events[0].Start
	var latest = events[0].End

	// We find out the edges for our time table
	for _, evt := range events {
		if evt.Start.Hour()-earliest.Hour() < 0 {
			earliest = evt.Start
		}
		if evt.Start.Hour()-earliest.Hour() > 0 {
			latest = evt.End
		}
	}

	empty := tableEvent{}
	// We prepare the table with what we know

	// We deal with the number of lines first :

	var howMuch = int(float64(60.0*(latest.Hour()-earliest.Hour())) / space.Minutes())
	for i := 0; i < howMuch; i++ {
		var line []tableEvent
		line = append(line, tableEvent{
			Start: earliest.Add(space * time.Duration(i)),
			Course: tableCourse{
				Name: "time",
			},
		})
		lines = append(lines, line)
	}
	// Then the days
	for i := range lines {
		for j := 0; j < lastDay-firstDay+1; j++ {
			lines[i] = append(lines[i], empty)
		}
	}

	// We then handle the course that might not start on the space timing
	for _, evt := range events {
		if timeToMin(evt.Start)%int(space.Minutes()) != 0 {
			// We need to add a row for the beginning
			for i, line := range lines {
				if line[0].Start.Hour()-evt.Start.Hour() > 0 {
					lines = append(lines[:i], lines[i-1:]...)
					var current []tableEvent
					current = append(current, tableEvent{
						Start: evt.Start,
						Course: tableCourse{
							Name: "time",
						},
					})
					for j := 0; j < lastDay-firstDay+1; j++ {
						current = append(current, empty)
					}
					lines[i-1] = current
					break
				}
			}
		}
		if timeToMin(evt.Start)%int(space.Minutes()) != 0 {
			// We need to add a row for the beginning
			for i, line := range lines {
				if line[0].Start.Hour()-evt.End.Hour() > 0 {
					lines = append(lines[:i], lines[i-1:]...)
					var current []tableEvent
					current = append(current, tableEvent{
						Start: evt.End,
						Course: tableCourse{
							Name: "time",
						},
					})
					for j := 0; j < lastDay-firstDay+1; j++ {
						current = append(current, empty)
					}
					lines[i] = current
					break
				}
			}
		}
	}

	// Now we fill in the blanks
	for i, line := range lines {
		for j := range line {
			if j == 0 {
				// This displays the time
				continue
			}
			origin := line[0].Start
			for _, evt := range events {
				if j-1 != evt.Start.Day()-firstDay {
					continue
				}
				if timeToMin(evt.Start) == timeToMin(origin) {
					lines[i][j] = evt
				}
				if i == 0 {
					continue
				}
				if timeToMin(origin) < timeToMin(evt.End) && timeToMin(lines[i-1][j].Start) >= timeToMin(evt.Start) {
					lines[i][j] = evt
				}
			}
		}
	}

	// SEND IT

	resp, _ := client.Database.Get(context.Background(), notionapi.DatabaseID(db))

	client.Database.Update(context.Background(), notionapi.DatabaseID(db), &notionapi.DatabaseUpdateRequest{
		Properties: resp.Properties,
	})

	for i := len(lines) - 1; i > -1; i-- {
		line := lines[i]
		properties := make(map[string]notionapi.Property)
		propertiesNoColor := make(map[string]notionapi.Property)
		for j, evt := range line {
			if j == 0 {
				properties[column[j]] = notionapi.TitleProperty{
					Title: []notionapi.RichText{
						{
							Text: notionapi.Text{
								Content: fmt.Sprintf("%02dh%02d", evt.Start.Hour(), evt.Start.Minute()),
							},
						},
					},
				}
				propertiesNoColor[column[j]] = notionapi.TitleProperty{
					Title: []notionapi.RichText{
						{
							Text: notionapi.Text{
								Content: fmt.Sprintf("%02dh%02d", evt.Start.Hour(), evt.Start.Minute()),
							},
						},
					},
				}
			} else {
				if evt.Course.Name == "" {
					continue
				}
				color := notionapi.ColorGray
				if strings.Contains(evt.Course.Name, "TD") {
					color = notionapi.ColorBrown
				} else if strings.Contains(evt.Course.Name, "CM") {
					color = notionapi.ColorBlue
				} else if strings.Contains(evt.Course.Name, "TEST") {
					color = notionapi.ColorGreen
				} else if strings.Contains(evt.Course.Name, "Evaluation") {
					color = notionapi.ColorOrange
				}
				properties[column[j]] = notionapi.SelectProperty{
					Select: notionapi.Option{
						Name:  evt.Course.Name,
						Color: color,
					},
				}
				propertiesNoColor[column[j]] = notionapi.SelectProperty{
					Select: notionapi.Option{
						Name: evt.Course.Name,
					},
				}
			}
		}

		_, err := client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
			Parent: notionapi.Parent{
				DatabaseID: notionapi.DatabaseID(resp.ID),
			},
			Properties: properties,
		})
		if err != nil {
			client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
				Parent: notionapi.Parent{
					DatabaseID: notionapi.DatabaseID(resp.ID),
				},
				Properties: propertiesNoColor,
			})
		}
	}
}
