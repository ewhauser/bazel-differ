package internal

import (
	"bytes"
	"os"
	"testing"
)

func TestDelimitedReader(t *testing.T) {
	protoBytes, err := os.ReadFile("query.protodelim")
	if err != nil {
		t.Errorf("Error reading file")
	}
	reader := NewReader(bytes.NewReader(protoBytes))
	targets, err := reader.ReadTargets()
	if err != nil {
		t.Errorf("error marshalling: %s", err)
		t.FailNow()
	}
	if len(targets) != 49 {
		t.Errorf("Expecting 49 targets but got %d", len(targets))
	}
}
