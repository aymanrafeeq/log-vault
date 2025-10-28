package parser

import (
	"fmt"
	"logGen/model"
	"regexp"
	"time"
)

var logPattern = regexp.MustCompile(`^(?P<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+\|\s+(?P<level>[A-Z]+)\s+\|\s+(?P<component>[\w-]+)\s+\|\s+host=(?P<host>[\w-]+)\s+\|\s+request_id=(?P<request_id>[\w-]+)\s+\|\s+msg="(?P<msg>.*)"$`)

func ParseLogEntry(line string) (*model.LogEntry, error) {
	match := logPattern.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("invalid format")
	}
	result := make(map[string]string)
	for i, name := range logPattern.SubexpNames() {
		if name != "" {
			result[name] = match[i]
		}
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05.000", result["time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}
	entry := model.LogEntry{
		Time:      parsedTime,
		Level:     result["level"],
		Component: result["component"],
		Host:      result["host"],
		ReqID:     result["request_id"],
		Msg:       result["msg"],
		Raw:       match[0],
	}
	return &entry, nil
}
