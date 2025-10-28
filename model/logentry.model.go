package model

import "time"

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
	ERROR LogLevel = "ERROR"
)

type LogEntry struct {
	Time      time.Time
	Level     string
	Component string
	Host      string
	ReqID     string
	Msg       string
	Raw       string
}
