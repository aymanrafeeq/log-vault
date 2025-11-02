package main

import (
	"flag"
	"fmt"
	"log/slog"
	"logGen/filter"
	"logGen/segment"
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

func main() {

	segments, err := segment.ParseLogSegments("/home/ayman/log-vault/logs")

	if err != nil {
		slog.Error("Error in creating segments: ", "Error", err)
	}

	level := flag.String("level", "", "Filter by log level")
	component := flag.String("component", "", "Filter by component")
	host := flag.String("host", "", "Filter by host")
	reqID := flag.String("reqID", "", "Filter by requestID")
	startTimeString := flag.String("after", "", "Filter by start time")
	endTimeString := flag.String("before", "", "Filter by end time")
	flag.Parse()

	split := func(s string) []string {
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}

	levels := split(*level)
	components := split(*component)
	hosts := split(*host)
	reqIDs := split(*reqID)

	var startTime, endTime time.Time

	if *startTimeString != "" {
		startTime, err = time.Parse("2006-01-02 15:04:05", *startTimeString)
		if err != nil {
			slog.Error("Error parsing start time", "error", err)
		}
	}

	if *endTimeString != "" {
		endTime, err = time.Parse("2006-01-02 15:04:05", *endTimeString)
		if err != nil {
			slog.Error("Error parsing end time", "error", err)
		}
	}

	filteredLogs := filter.FilterEntries(segments, levels, components, hosts, reqIDs, startTime, endTime)

	for _, entry := range filteredLogs {
		fmt.Println(entry.Raw)
	}
	fmt.Printf("Found %d matching entries\n", len(filteredLogs))
}
