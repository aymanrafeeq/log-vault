package segment

import (
	"logGen/model"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func createTestFile(t *testing.T, dir string, name string, entry []string) (string, error) {
	t.Helper()
	path := filepath.Join(dir, name)
	content := strings.Join(entry, "\n")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v\n", err)
	}
	return dir, nil
}

func TestParseLogSegments(t *testing.T) {
	temDir := t.TempDir()
	name := "file1"

	data := []string{
		`2025-10-23 16:53:00.033 | WARN | auth | host=worker01 | request_id=req-123 | msg="Cache server connected"`,
		`2025-10-23 16:54:00.000 | INFO | api | host=web01 | request_id=req-234 | msg="API call processed"`,
	}
	path, _ := createTestFile(t, temDir, name, data)

	got, err := ParseLogSegments(path)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(got) == 0 {
		t.Fatalf("Expected at least one segment, got none")
	}

	expectedStart, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-23 16:53:00.033")
	expectedEnd, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-23 16:54:00.000")

	if got[0].StartTime != expectedStart {
		t.Errorf("Expected StartTime %v, got %v", expectedStart, got[0].StartTime)
	}
	if got[0].EndTime != expectedEnd {
		t.Errorf("Expected EndTime %v, got %v", expectedEnd, got[0].EndTime)
	}

	expectedIndex := model.SegmentIndex{
		ByLevel: map[string][]int{
			"WARN": {0},
			"INFO": {1},
		},
		ByComponent: map[string][]int{
			"auth": {0},
			"api":  {1},
		},
	}

	if !reflect.DeepEqual(got[0].Index.ByLevel, expectedIndex.ByLevel) {
		t.Errorf("Expected ByLevel: %v, got: %v", expectedIndex.ByLevel, got[0].Index.ByLevel)
	}
}

func TestParseLogSegments_InvalidPath(t *testing.T) {

	_, err := ParseLogSegments("bhuhu")
	if err == nil {
		t.Errorf("Expected error for invalid path, got nil")
	}
}
