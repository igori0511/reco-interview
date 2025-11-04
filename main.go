package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/time/rate"
)

const (
	asanaBaseURL = "https://app.asana.com/api/1.0"
	workspaceGid = "1211834271836941"
	// !! WARNING: Do not hardcode tokens in production code.
	// This should be loaded from environment variables or a secure vault.
	bearerToken = "2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"
	// --- New Configuration for Saving ---
	dataFolderName = "extracted-asana-data"
	dataFileName   = "asana_data.json"

	maxRetries = 3 // Maximum number of retries for a single request
)

var extractFn = extractAsanaUsers

// Asana's standard plan limit is 150 requests/minute.
var asanaLimiter = rate.NewLimiter(rate.Every(time.Minute/150), 10)

func main() {
	intervalSec := flag.Int("interval", 5, "allowed: 5 (seconds) or 1800 (30 minutes)")
	flag.Parse()

	if *intervalSec != 5 && *intervalSec != 1800 {
		log.Fatalf("invalid -interval=%d (allowed: 5 or 1800)", *intervalSec)
	}
	interval := time.Duration(*intervalSec) * time.Second

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// run once immediately
	if err := extractFn(ctx); err != nil {
		log.Printf("extractAsanaUsers (initial): %v", err)
	}

	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case <-timer.C:
			if err := extractFn(ctx); err != nil {
				log.Printf("extractAsanaUsers: %v", err)
			}
		}
	}
}

// TODO structure it a bit better since we can use interfaces and structs.
func extractAsanaUsers(ctx context.Context) error {
	associatedUsers, err := getAssociatedAsanaUsers(ctx)
	if err != nil {
		return err
	}
	err = saveDataToFile(associatedUsers)
	if err != nil {
		return err
	}
	return nil
}
