package main

import (
	"fmt"
	"log"
	models "logGen/pkg/dbmodels"
	"logGen/pkg/parser"
	"logGen/pkg/web"
	"os"
)

const dbUrl = "postgresql:///logVault?host=/var/run/postgresql/"

func handleCommand(args []string) error {
	db, err := models.CreateDB(dbUrl)
	if err != nil {
		return err
	}
	switch args[0] {
	case "init":
		err := models.InitDb(db)
		if err != nil {
			return err
		}
	case "add":
		folderName := args[1]

		entries, err := parser.ParseLogFiles(folderName)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}

		for _, entry := range entries {
			models.AddEntry(db, entry)
		}

		return nil
	case "query":
		queryList := args[1:]

		entries, err := models.Query(db, queryList)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fmt.Println(entry)
		}
		fmt.Printf("%d entries matched: \n", len(entries))
		return nil

	case "web":
		r := web.SetupRoutes(db)
		log.Println("Server running at http://localhost:8080")
		return r.Run(":8080")

	default:
		return fmt.Errorf("unknown command: %s (expected: init | add | query | web)", args[0])
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		// default to running web if no args (optional)
		fmt.Fprintf(os.Stderr, "Usage: %s <init|add|query|web> ...\n", os.Args[0])
		os.Exit(2)
	}
	err := handleCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in invocation %v\n", err)
		os.Exit(-1)
	}
}
