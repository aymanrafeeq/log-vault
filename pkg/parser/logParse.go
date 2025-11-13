package parser

import (
	"bufio"
	"fmt"
	"log/slog"
	"logGen/model"
	"os"
	"path/filepath"
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

func ParseLogFiles() ([]model.LogEntry, error) {
	var allEntries []model.LogEntry

	folderPath := "/home/ayman/log-vault/logs"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory : %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(folderPath, file.Name())
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v\n", path, err)
			continue

		}
		defer f.Close()
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Text()
			entry, err := ParseLogEntry(line)
			if err != nil {
				slog.Error("Error while parsing : ", "error", err)
				continue
			}
			allEntries = append(allEntries, *entry)
		}

	}
	return allEntries, nil
}
