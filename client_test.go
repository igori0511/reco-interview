package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

type asanaProjectsResp struct {
	Data []struct {
		Gid          string `json:"gid"`
		Name         string `json:"name"`
		ResourceType string `json:"resource_type"`
	} `json:"data"`
}

func parseProjectsResp(body []byte) (*asanaProjectsResp, error) {
	var u asanaProjectsResp
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func TestMakeAsanaRequest_Success200(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/workspaces/%s/projects", asanaBaseURL, workspaceGid)
	body, err := makeAsanaRequest(ctx, url)
	if err != nil {
		t.Fatalf("makeAsanaRequest error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("empty body from Asana")
	}
	// Basic sanity: Asana JSON responses include "data"
	if !bytes.Contains(body, []byte(`"data"`)) {
		t.Fatalf("unexpected response (no \"data\"): %s", string(body))
	}

	// Parse and validate expected fields
	resp, err := parseProjectsResp(body)
	if resp == nil || len(resp.Data) == 0 {
		t.Fatalf("unexpected response (no \"data\"): %s", string(body))
	}
	first := resp.Data[0]
	if first.Gid == "" || first.Name == "" {
		t.Fatalf("project missing gid or name; body: %s", string(body))
	}
	if first.ResourceType != "" && first.ResourceType != "project" {
		t.Fatalf("unexpected resource_type: %q (want \"project\")", first.ResourceType)
	}

	// Optional strict checks if you know the exact project
	if wantGID := os.Getenv("ASANA_EXPECT_PROJECT_GID"); wantGID != "" && first.Gid != wantGID {
		t.Fatalf("first project gid = %q, want %q", first.Gid, wantGID)
	}
	if wantName := os.Getenv("ASANA_EXPECT_PROJECT_NAME"); wantName != "" && first.Name != wantName {
		t.Fatalf("first project name = %q, want %q", first.Name, wantName)
	}
}
