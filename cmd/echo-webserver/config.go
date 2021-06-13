package main

import (
	"encoding/json"
	"log"
	"os"
)

type Provider struct {
	ClientID     string
	ClientSecret string
}

type Configuration struct {
	GoogleOAuth   Provider
	FacebookOAuth Provider
	SessionsKey   string
	DSN           string
	Port          string
}

var (
	GlobalConfig = new(Configuration)
)

func LoadConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(GlobalConfig); err != nil {
		log.Fatal(err)
	}
}
