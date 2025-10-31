package main

import (
	"fmt"
	"logGen/segment"
)

func main() {
	// line := `2025-10-26 11:16:12.840 | DEBUG | cache | host=web01 | request_id=req-ymuon4-1921 | msg="Connection established to replica"`

	// entry, err := parser.ParseLogEntry(line)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("struct: %#v", entry)

	// 	entries, err := parser.ParseLogFiles()
	// 	if err != nil {
	// 		slog.Error("Error here", "error", err)
	// 	}
	// 	fmt.Println(entrpies)
	segments, _ := segment.ParseLogSegments("/home/ayman/log-vault/logs")

	fmt.Printf("filename: %#v \n", segments)
	// fmt.Println("filename: %#v \n", segments[1].FileName)
	// for _, segment := range segments {
	// 	fmt.Printf("File Name : %s\n", segment.FileName)
	// 	fmt.Printf("Start Time : %v\n", segment.StartTime)
	// 	fmt.Printf("End Time : %v\n\n\n\n", segment.EndTime)
	// 	// fmt.Printf("LogEntries : %v", segment.LogEntries)
	// }

}
