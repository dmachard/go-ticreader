package ticreader

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseHistoricFrame(t *testing.T) {
	frame := "ADCO 012345678901 C"
	group := parseHistoricFrame(frame)

	if group.Label != "ADCO" {
		t.Errorf("Expected label ADCO, got %s", group.Label)
	}

	if group.Data != "012345678901" {
		t.Errorf("Expected data 012345678901, got %s", group.Data)
	}
}

func TestParseStandardFrame(t *testing.T) {
	frame := "OPTARIF\tBASE\tC"
	group := parseStandardFrame(frame)

	if group.Label != "OPTARIF" {
		t.Errorf("Expected label OPTARIF, got %s", group.Label)
	}

	if group.Data != "BASE" {
		t.Errorf("Expected data BASE, got %s", group.Data)
	}
}

func TestReadFrame(t *testing.T) {
	r := strings.NewReader("\x02ADCO 012345678901 C\x03")
	reader := bufio.NewReader(r)

	frame, err := readFrame(reader, ModeHistorical)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if frame.Timestamp.IsZero() {
		t.Errorf("Timestamp should not be zero")
	}
}
