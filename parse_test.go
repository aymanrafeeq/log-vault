package main

import (
	"testing"
	"time"
)

func TestParseLogEntry(t *testing.T) {
	line := `2025-10-26 11:16:12.840 | DEBUG | cache | host=web01 | request_id=req-ymuon4-1921 | msg="Connection established to replica"`
	entry, err := ParseLogEntry(line)

	if err != nil {
		t.Errorf("Log Parsing Failed! got error: %v", err)
	}

	expectedTime, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-26 11:16:12.840")

	if entry.Raw != line {
		t.Errorf("Expected raw to be %q but got %q", line, entry.Raw)
	}

	if !entry.Time.Equal(expectedTime) {
		t.Errorf("Expected time %v but got %v", expectedTime, entry.Time)
	}

	if entry.Level != "DEBUG" {
		t.Errorf("Expected DEBUG but got %s\n", entry.Level)
	}

	if entry.Component != "cache" {
		t.Errorf("Expected cache but got %s\n", entry.Component)
	}

	if entry.Host != "web01" {
		t.Errorf("Expected web01 but got %s\n", entry.Host)
	}

	if entry.ReqID != "req-ymuon4-1921" {
		t.Errorf("Expected req-ymuon4-1921 but got %s\n", entry.ReqID)
	}

	if entry.Msg != "Connection established to replica" {
		t.Errorf("Expected 'Connection established to replica' but got %s\n", entry.Msg)
	}
}
