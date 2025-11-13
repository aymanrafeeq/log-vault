package main

import (
	"context"
	"flag"
	"log/slog"
	"logGen/model"
	"logGen/pkg/segment"
	"os"

	"github.com/jackc/pgx/v5"
)

func InsertLogs(ctx context.Context, conn *pgx.Conn, segments []model.Segment) error {
	for _, segment := range segments {
		for _, entry := range segment.LogEntries {
			_, err := conn.Exec(ctx, `
                INSERT INTO logEntry (timestamp, level, host, component, reqID, message)
                VALUES (
                    $1,
                    (SELECT id FROM logLevel WHERE level = $2),
                    (SELECT id FROM logHost WHERE host = $3),
                    (SELECT id FROM logComponent WHERE component = $4),
                    $5,
                    $6
                )`,
				entry.Time,
				entry.Level,
				entry.Host,
				entry.Component,
				entry.ReqID,
				entry.Msg,
			)
			if err != nil {
				slog.Warn("Failed to insert log entry", "error", err, "entry", entry.Raw)
				continue
			}
		}
	}
	return nil
}
func main() {
	logPath := flag.String("path", "/home/ayman/log-vault/logs", "Path to the log directory")
	flag.Parse()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	segments, err := segment.ParseLogSegments(*logPath)
	if err != nil {
		slog.Error("Failed to parse logs", "error", err)
	}
	err = InsertLogs(ctx, conn, segments)
	if err != nil {
		slog.Error("Failed to insert logs", "error", err)
	} else {
		slog.Info("All logs inserted successfully")
	}
}
