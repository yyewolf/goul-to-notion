package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var refreshToken = os.Getenv("initial_refresh")
var token = os.Getenv("initial_token")
var notionSecret = os.Getenv("notion_secret")
var db = os.Getenv("notion_db")

var file = "/data/var.json"

type Config struct {
	Token        string `json:"t"`
	RefreshToken string `json:"r"`
	Secret       string `json:"s"`
	DB           string `json:"d"`
}

func saveVar() {
	conf := &Config{
		Token:        token,
		RefreshToken: refreshToken,
		Secret:       notionSecret,
		DB:           db,
	}
	d, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("err saving to file.")
		return
	}
	err = os.WriteFile(file, d, os.ModeAppend)
	if err != nil {
		fmt.Println("err saving to file.")
		return
	}
}

func getVar() {
	d, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("couldn't read file.")
		saveVar()
		return
	}
	conf := &Config{}
	json.Unmarshal(d, conf)
	fmt.Println("read saved file.")
}
