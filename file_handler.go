package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// getAssociatedAsanaUsers performs the full logic:
// 1. Fetches all projects.
// 2. Loops through each project and fetches its members.
// 3. Returns a map of [Project Name] -> [List of Member Names]
func getAssociatedAsanaUsers() (map[string][]string, error) {
	// This is the final data structure we will return.
	// It maps a project name (string) to a list of member names ([]string).
	projectUserMap := make(map[string][]string)

	// Step 1: Get all projects in the workspace.
	projectsResp, err := getProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	fmt.Printf("Found %d projects. Now fetching members for each...\n", len(projectsResp.Data))

	// Step 2: Loop through each project.
	for _, project := range projectsResp.Data {
		var memberNames []string

		// Step 3: Get the members for this specific project.
		membershipsResp, err := getProjectMembers(project.GID)
		if err != nil {
			// Log the error for this project but continue with others.
			fmt.Printf("Failed to get members for project %s (%s): %v\n", project.Name, project.GID, err)
			continue // Skip to the next project
		}

		// Extract just the names for our final map.
		for _, membership := range membershipsResp.Data {
			memberNames = append(memberNames, membership.Member.Name)
		}

		// Add the results to our map.
		fmt.Printf("  - Project '%s' has %d members: %v\n", project.Name, len(memberNames), memberNames)
		projectUserMap[project.Name] = memberNames
	}

	return projectUserMap, nil
}

func saveDataToFile(data map[string][]string) error {
	fmt.Printf("\n--- Saving data to file ---\n")

	// 1. Create the folder if it doesn't exist
	// 0755 provides standard directory permissions (read/write/execute for owner, read/execute for others)
	if err := os.MkdirAll(dataFolderName, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dataFolderName, err)
	}

	// 2. Define the full file path
	filePath := filepath.Join(dataFolderName, dataFileName)

	// 3. Marshal the data map into pretty-formatted JSON
	jsonData, err := json.MarshalIndent(data, "", "  ") // "  " for 2-space indentation
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// 4. Write the JSON data to the file
	// 0644 provides standard file permissions (read/write for owner, read for others)
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write data to file %s: %w", filePath, err)
	}

	fmt.Printf("Successfully saved data to: %s\n", filePath)
	return nil
}
