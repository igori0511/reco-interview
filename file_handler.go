package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// getAssociatedAsanaUsers performs the full logic:
// 1. Fetches all projects.
// 2. Loops through each project and fetches its members.
// 3. Returns a map of [Project Name] -> [List of Member Names]
func getAssociatedAsanaUsers(ctx context.Context) (map[string][]string, error) {
	projectUserMap := make(map[string][]string)

	projectsResp, err := getProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	for _, project := range projectsResp.Data {
		var memberNames []string
		membershipsResp, err := getProjectMembers(ctx, project.GID)
		if err != nil {
			fmt.Printf("Failed to get members for project %s (%s): %v\n", project.Name, project.GID, err)
			continue // Skip to the next project
		}

		for _, membership := range membershipsResp.Data {
			memberNames = append(memberNames, membership.Member.Name)
		}

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
	filePath := filepath.Join(dataFolderName, dataFileName)

	jsonData, err := json.MarshalIndent(data, "", "  ") // "  " for 2-space indentation
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Write the JSON data to the file
	// 0644 provides standard file permissions (read/write for owner, read for others)
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write data to file %s: %w", filePath, err)
	}

	fmt.Printf("Successfully saved data to: %s\n", filePath)
	return nil
}
