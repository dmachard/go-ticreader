package ticreader

import (
	"bufio"
	"strings"
	"testing"
)

func TestDecodeFrame(t *testing.T) {
	r := strings.NewReader("\x020x0AADCO 012345678901 E0x0D\x03")
	reader := bufio.NewReader(r)

	decodedFrame, err := decodeFrame(reader, ModeHistorical)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if decodedFrame.Timestamp.IsZero() {
		t.Errorf("Timestamp should not be zero")
	}
}

func TestDecodeDataset(t *testing.T) {
	frame := "\x0AADCO 123445678901 H\x0D\x0AOPTARIF BASE 0\x0D\x0AISOUSC 45 ?\x0D"
	decodedFrame := decodeDataset(frame, ModeHistorical)

	if len(decodedFrame.Dataset) != 3 {
		t.Errorf("Expected 3 groups, got %v", len(decodedFrame.Dataset))
	}

	if decodedFrame.Dataset[0].Label != "ADCO" {
		t.Errorf("Expected first label to be ADCO, got %s", decodedFrame.Dataset[0].Label)
	}
}

func TestParseHistoricDataset(t *testing.T) {
	frame := "ADCO 012345678901 E"
	dataset, err := parseHistoricDataset(frame)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if dataset.Label != "ADCO" {
		t.Errorf("Expected label ADCO, got %s", dataset.Label)
	}

	if dataset.Data != "012345678901" {
		t.Errorf("Expected data 012345678901, got %s", dataset.Data)
	}

	if dataset.Valid != true {
		t.Errorf("Checksum invalid: %s", dataset.Checksum)
	}
}

func TestParseStandardDataset(t *testing.T) {
	frame := "ADCO\t1234567890123\t,"
	dataset, err := parseStandardDataset(frame)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if dataset.Label != "ADCO" {
		t.Errorf("Expected label ADCO, got %s", dataset.Label)
	}

	if dataset.Data != "1234567890123" {
		t.Errorf("Expected data 1234567890123, got %s", dataset.Data)
	}

	if dataset.Valid != true {
		t.Errorf("Checksum invalid, got %s", dataset.Checksum)
	}
}

func TestCalculateChecksum(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"IINST 002", "Y"}, // historical format
		{"IMAX 090", "H"},
		{"SINSTS\t00520\t", "M"},                // standard format sans horodatage
		{"SMAXSN\tH250221001022\t02560\t", "+"}, // avec horodatages
		{"STGE\t003AC001\t", "M"},
		{"ADCO\t1234567890123\t", ","},
	}

	for _, test := range tests {
		result := calculateChecksum(test.input)
		if string(result) != test.expected {
			t.Errorf("calculateChecksum(%q) = %s; expected %s", test.input, string(result), test.expected)
		}
	}
}
