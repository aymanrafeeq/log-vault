package segment

import (
	"bufio"
	"fmt"
	"logGen/model"
	"logGen/pkg/parser"
	"os"
	"path/filepath"
	"time"
)

func ParseLogSegments(path string) ([]model.Segment, error) {
	var segments []model.Segment

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(path, file.Name())
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		var entries []model.LogEntry
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			entry, err := parser.ParseLogEntry(line)
			if err == nil && entry != nil {
				entries = append(entries, *entry)
			}
		}

		if len(entries) == 0 {
			continue
		}

		segment := model.Segment{
			FileName:   file.Name(),
			LogEntries: entries,
			//StartTime:  entries[0].Time,
			//EndTime:    entries[len(entries)-1].Time,
			Index: BuildSegmentIndex(entries),
		}
		for _, entry := range entries {
			if segment.StartTime.Equal(time.Time{}) || entry.Time.Before(segment.StartTime) {
				segment.StartTime = entry.Time
			}

			if entry.Time.After(segment.EndTime) {
				segment.EndTime = entry.Time
			}
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

func BuildSegmentIndex(LogEntries []model.LogEntry) model.SegmentIndex {
	index := model.SegmentIndex{
		ByLevel:     make(map[string][]int),
		ByComponent: make(map[string][]int),
		ByHost:      make(map[string][]int),
		ByReqId:     make(map[string][]int),
	}
	for idx, LogEntry := range LogEntries {
		index.ByLevel[string(LogEntry.Level)] = append(index.ByLevel[string(LogEntry.Level)], idx)
		index.ByComponent[LogEntry.Component] = append(index.ByComponent[LogEntry.Component], idx)
		index.ByHost[LogEntry.Host] = append(index.ByHost[LogEntry.Host], idx)
		index.ByReqId[LogEntry.ReqID] = append(index.ByReqId[LogEntry.ReqID], idx)
	}
	return index
}
