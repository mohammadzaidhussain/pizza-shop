package logger

import (
	"log"
	"os"
)

func Log(message any) {
	isLogenabled := os.Getenv("log")
	if isLogenabled != "" {
		log.Println(message)
	}
}
