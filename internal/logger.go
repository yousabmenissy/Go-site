package internal

import (
	"log"
	"os"
)

const (
	bold               = "\033[1m"
	reset              = "\033[0m"
	green              = "\033[38;2;85;255;127m"
	purple             = "\033[38;2;207;152;244m"
	red                = "\033[31m"
	yeallow            = "\033[38;2;255;246;151m"
	boldGreen_INFO     = bold + green + "  INFO    " + reset
	boldPurple_DEBUG   = bold + purple + "  DEBUG   " + reset
	boldRed_Error      = bold + red + "  Error   " + reset
	boldYellow_WARNING = bold + yeallow + "  WARNING " + reset
)

type ConsoleLogger struct {
	LogInfo    *log.Logger
	LogDebug   *log.Logger
	LogWarning *log.Logger
	LogError   *log.Logger
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		LogInfo:    log.New(os.Stdout, boldGreen_INFO, log.Ldate|log.Ltime|log.Lshortfile),
		LogDebug:   log.New(os.Stdout, boldPurple_DEBUG, log.Ldate|log.Ltime|log.Lshortfile),
		LogWarning: log.New(os.Stdout, boldYellow_WARNING, log.Ldate|log.Ltime|log.Lshortfile),
		LogError:   log.New(os.Stderr, boldRed_Error, log.Ldate|log.Ltime|log.Lshortfile),
	}
}
