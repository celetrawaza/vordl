package main

import (
	"flag"
	"time"
)

type ConfigStruct struct {
	Port            int
	ResetInterval   time.Duration
	StartWordNumber int
}

var Config ConfigStruct

func LoadConfig() ConfigStruct {
	port := flag.Int("port", 8080, "Port to serve the app")
	wordNumber := flag.Int("wordNumber", 1, "Start word number")
	flag.Parse()
	return ConfigStruct{
		Port:            *port,
		ResetInterval:   time.Hour,
		StartWordNumber: *wordNumber,
	}
}
