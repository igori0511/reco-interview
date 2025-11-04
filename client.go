package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// ProjectsResponse maps to the JSON from GET /workspaces/{...}/projects
type ProjectsResponse struct {
	Data []struct {
		GID  string `json:"gid"`
		Name string `json:"name"`
	} `json:"data"`
}

// MembershipsResponse maps to the JSON from GET /memberships
type MembershipsResponse struct {
	Data []struct {
		Member struct {
			GID  string `json:"gid"`
			Name string `json:"name"`
		} `json:"member"`
	} `json:"data"`
}

func makeAsanaRequest(ctx context.Context, url string) ([]byte, error) {
	// Create a new HTTP client with a 10-second timeout.
	client := &http.Client{Timeout: 10 * time.Second}

	//TODO add exponential backoff functionality for better retry handling
	for i := 0; i < maxRetries; i++ {
		// Wait for the token bucket limiter to allow this request.
		// This will block if we are sending requests too quickly.
		if err := asanaLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}

		// Create a new GET request.
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create new request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+bearerToken)
		req.Header.Set("Accept", "application/json")

		// Send the request.
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		switch resp.StatusCode {
		case http.StatusOK: // 200 OK
			// Success! Read the body and return.
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %w", err)
			}
			return body, nil

		case http.StatusTooManyRequests: // 429 Rate Limit
			retryAfterHeader := resp.Header.Get("Retry-After")
			resp.Body.Close()

			retryAfterSeconds, parseErr := strconv.Atoi(retryAfterHeader)
			if parseErr != nil {
				// If parsing fails, Asana docs mention a 1-minute "cost".
				// We will default to a safe 60-second wait.
				fmt.Printf("Rate limit hit. 'Retry-After' header unreadable ('%s'). Waiting 60 seconds...\n", retryAfterHeader)
				retryAfterSeconds = 60
			}

			// Add a small buffer to be safe
			waitDuration := time.Duration(retryAfterSeconds+1) * time.Second

			fmt.Printf("Rate limit hit for request: %s\n", url)
			fmt.Printf("Retrying attempt %d/%d after %s...\n", i+1, maxRetries, waitDuration)

			// Sleep for the required duration
			time.Sleep(waitDuration)

		default:
			// Any other error (401, 403, 404, 500, etc.) is a hard failure.
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("asana API error: status %d, body: %s", resp.StatusCode, string(body))
		}
	}

	// If we exit the loop, it means we exceeded maxRetries.
	return nil, fmt.Errorf("failed to make request after %d retries: %s", maxRetries, url)
}

// getProjects fetches all projects for the configured workspace.
func getProjects(ctx context.Context) (ProjectsResponse, error) {
	url := fmt.Sprintf("%s/workspaces/%s/projects", asanaBaseURL, workspaceGid)
	var projects ProjectsResponse

	body, err := makeAsanaRequest(ctx, url)
	if err != nil {
		return projects, err
	}

	// Unmarshal the JSON response into our struct.
	if err := json.Unmarshal(body, &projects); err != nil {
		return projects, fmt.Errorf("failed to unmarshal projects: %w", err)
	}

	return projects, nil
}

/* TODO this can be done in parallel in respect to limits of asana api
 so the usage of multiple goroutines can be leveraged, for each project make a requests and limit the concurrency with channels.
this will speed up the process tremendously
*/
// getProjectMembers fetches all members for a single given project GID.
func getProjectMembers(ctx context.Context, projectGid string) (MembershipsResponse, error) {
	url := fmt.Sprintf("%s/memberships?parent=%s", asanaBaseURL, projectGid)
	var memberships MembershipsResponse

	body, err := makeAsanaRequest(ctx, url)
	if err != nil {
		return memberships, err
	}

	// Unmarshal the JSON response into our struct.
	if err := json.Unmarshal(body, &memberships); err != nil {
		return memberships, fmt.Errorf("failed to unmarshal memberships: %w", err)
	}

	return memberships, nil
}
