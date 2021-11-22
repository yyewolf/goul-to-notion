package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type tableResponse struct {
	Data tableData `json:"data"`
}

type tableData struct {
	Table tableTable `json:"timetable"`
}

type tableTable struct {
	Plannings []tablePlanning `json:"plannings"`
}

type tablePlanning struct {
	Events []tableEvent `json:"events"`
}

type tableEvent struct {
	Start  time.Time    `json:"startDateTime"`
	End    time.Time    `json:"endDateTime"`
	Course tableCourse  `json:"course"`
	Rooms  []tableRoom  `json:"rooms"`
	Groups []tableGroup `json:"groups"`
}

type tableCourse struct {
	Name  string `json:"label"`
	Color string `json:"color"`
}

type tableRoom struct {
	Name string `json:"label"`
}

type tableGroup struct {
	Name string `json:"label"`
}

func getTablePrimitive(token string, from, to int64) (response tableResponse) {
	body := fmt.Sprintf(`{
		"operationName":"timetable",
		"variables":{
		   "uid":"smagghe2u",
		   "from":%d000,
		   "to":%d000
		},
		"query":"query timetable($uid: String!, $from: Float, $to: Float) {\n  timetable(uid: $uid, from: $from, to: $to) {\n    plannings {\n      events {\n        startDateTime\n        endDateTime\n        course {\n          label\n          color\n        }\n        rooms {\n          label\n        }\n        groups {\n          label\n        }\n      }\n    }\n  }\n}"
	 }`, from, to)
	req, err := http.NewRequest("POST", "https://multi.univ-lorraine.fr/graphql", strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	req.Header.Add("x-refresh-token", refreshToken)
	req.Header.Add("x-token", token)
	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.Header.Get("X-Token-Status") == "201" {
		refreshToken = resp.Header.Get("X-Refresh-Token")
		token = resp.Header.Get("X-Token")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		panic(err)
	}
	return
}
