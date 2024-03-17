package helpers

import (
	"log"
)

const (
	Info    = "info"
	Error   = "error"
	Warning = "warning"
	Debug   = "debug"
)

func Log(level, message string) {
	switch level {
	case Info:
		log.Printf("\u001b[34mINFO: \u001B[0m %s", message)
	case Error:
		log.Printf("\u001b[31mERROR: \u001b[0m %s", message)
	case Warning:
		log.Printf("\u001b[33mWARNING: \u001B[0m %s", message)
	case Debug:
		log.Printf("\u001b[36mDEBUG: \u001B[0m %s", message)
	}
}
