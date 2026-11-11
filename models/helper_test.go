package models

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var testData *Data

func TestMain(m *testing.M) {
	var err error
	testData, err = NewData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load test data: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func jsonOf(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}
