package main

import (
	"fmt"
	"log/slog"
	"logGen/parser"
)

func main() {
	// line := `2025-10-26 11:16:12.840 | DEBUG | cache | host=web01 | request_id=req-ymuon4-1921 | msg="Connection established to replica"`

	// entry, err := parser.ParseLogEntry(line)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("struct: %#v", entry)

	entries, err := parser.ParseLogFiles()
	if err != nil {
		slog.Error("Error here", "error", err)
	}
	fmt.Println(entries)
}
