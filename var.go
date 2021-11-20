package main

import "os"

var refreshToken = os.Getenv("initial_refresh")
var token = os.Getenv("initial_token")
var notionSecret = os.Getenv("notion_secret")
var db = os.Getenv("notion_db")
