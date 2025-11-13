package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"logGen/filter"
	"logGen/model"
	"os"
	"strings"
	"time"
)

// func main() {
// 	// line := `2025-10-26 11:16:12.840 | DEBUG | cache | host=web01 | request_id=req-ymuon4-1921 | msg="Connection established to replica"`

// 	// entry, err := parser.ParseLogEntry(line)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Printf("struct: %#v", entry)

// 	// 	entries, err := parser.ParseLogFiles()
// 	// 	if err != nil {
// 	// 		slog.Error("Error here", "error", err)
// 	// 	}
// 	// 	fmt.Println(entrpies)
// 	segments, _ := segment.ParseLogSegments("/home/ayman/log-vault/logs")

// 	fmt.Printf("filename: %#v \n", segments)
// 	// fmt.Println("filename: %#v \n", segments[1].FileName)
// 	// for _, segment := range segments {
// 	// 	fmt.Printf("File Name : %s\n", segment.FileName)
// 	// 	fmt.Printf("Start Time : %v\n", segment.StartTime)
// 	// 	fmt.Printf("End Time : %v\n\n\n\n", segment.EndTime)
// 	// 	// fmt.Printf("LogEntries : %v", segment.LogEntries)
// 	// }

// }

func split(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func parseTime(value string, label string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	parsed, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		slog.Error("Error parsing time", "field", label, "input", value, "error", err)
		return time.Time{}, err
	}
	return parsed, nil
}

func main() {

	start := time.Now()
	jsonFile := flag.String("file", "logs.json", "Input JSON file generated")
	level := flag.String("level", "", "Filter by log level")
	component := flag.String("component", "", "Filter by component")
	host := flag.String("host", "", "Filter by host")
	reqID := flag.String("reqID", "", "Filter by requestID")
	startTimeString := flag.String("after", "", "Filter by start time")
	endTimeString := flag.String("before", "", "Filter by end time")
	flag.Parse()

	file, err := os.ReadFile(*jsonFile)

	if err != nil {
		slog.Error("Error in Opening JSON file: ", "Error", err)
		os.Exit(1)
	}

	var segments []model.Segment
	err = json.Unmarshal(file, &segments)

	if err != nil {
		slog.Error("Error in reading JSON file", "error", err)
	}

	levels := split(*level)
	components := split(*component)
	hosts := split(*host)
	reqIDs := split(*reqID)

	startTime, _ := parseTime(*startTimeString, "after")
	endTime, _ := parseTime(*endTimeString, "before")

	filteredLogs := filter.FilterEntries(segments, levels, components, hosts, reqIDs, startTime, endTime)

	for _, entry := range filteredLogs {
		fmt.Println(entry.Raw)
	}

	fmt.Printf("Found %d matching entries\n", len(filteredLogs))

	elapsed := time.Since(start)
	fmt.Printf("Total duration: %s\n", elapsed)
}
