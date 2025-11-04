package main

import (
	"log"
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

// Asana's standard plan limit is 150 requests/minute.
var asanaLimiter = rate.NewLimiter(rate.Every(time.Minute/150), 10)

func main() {

	// TODO
	/*
			Get projects
			1. curl "https://app.asana.com/api/1.0/workspaces/1211834271836941/projects" \
			     -H "Authorization: Bearer 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"

			For each project loop through:
			2. curl "https://app.asana.com/api/1.0/memberships?parent=1211834226422815" \
		     -H "Authorization: Bearer 2/1211834271836929/1211834383480110:098090c664d79859c951ba6858d531f3"

			I need a basic version that extract the data and saves it
			1. Http handler(basic without limits, expand later)
			2. Separate function that will create a map users -> projects, ot maybe just do in one function??? lets see
			3. Function that will save the data in json format in a file
			4. Write them into a file in json format
	*/

	associatedUsers, err := getAssociatedAsanaUsers()
	if err != nil {
		log.Fatal(err)
	}
	err = saveDataToFile(associatedUsers)
	if err != nil {
		log.Fatal(err)
	}
}
