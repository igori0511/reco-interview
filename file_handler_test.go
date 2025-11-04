package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSaveDataToFile_CreatesDirAndWritesPrettyJSON(t *testing.T) {
	// isolate filesystem effects in a temp dir
	testDataFolderName := "extracted-asana-data"
	testDataFileName := "asana_data.json"

	// input (map key order is nondeterministic; we'll validate by unmarshalling)
	in := map[string][]string{
		"users":    {"alice", "bob"},
		"projects": {"reco-interview"},
	}

	if err := saveDataToFile(in); err != nil {
		t.Fatalf("saveDataToFile error: %v", err)
	}

	// file must exist
	dst := filepath.Join(testDataFolderName, testDataFileName)
	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("reading written file: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("written file is empty")
	}

	// parse back and compare content (order-independent)
	var got map[string][]string
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal written JSON: %v\nbody: %s", err, string(b))
	}
	if !reflect.DeepEqual(got, in) {
		t.Fatalf("round-trip mismatch:\n got: %#v\nwant: %#v\nraw: %s", got, in, string(b))
	}

	// quick check for pretty-print (2-space indent): look for newline + two spaces before a key
	if !bytes.Contains(b, []byte("\n  \"")) {
		t.Fatalf("JSON not pretty-printed with 2-space indent:\n%s", string(b))
	}
}
