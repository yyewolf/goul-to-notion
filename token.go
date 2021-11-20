package main

import (
	"net/http"
	"strings"
	"time"
)

func doLoop() {
	t := time.NewTicker(15 * time.Minute)
	go func() {
		for {
			select {
			case <-t.C:
				smallRequest()
			}
		}
	}()
}

func smallRequest() {
	body := `{"operationName":"weather","variables":{},"query":"query weather {\n  weather {\n    name\n    main {\n      temp\n      __typename\n    }\n    weather {\n      description\n      id\n      __typename\n    }\n    __typename\n  }\n}\n"}`

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
	if resp.Header.Get("X-Token-Status") != "201" {
		return
	}
	refreshToken = resp.Header.Get("X-Refresh-Token")
	token = resp.Header.Get("X-Token")
}
