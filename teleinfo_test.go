package ticreader

import (
	"testing"

	"go.bug.st/serial"
)

func TestDecodeFrame(t *testing.T) {
	frame := "ADCO 012345678901 C\nOPTARIF BASE C\nISOUSC 30 9"
	result := decodeFrame(frame, ModeHistorical)

	if len(result.Informations) != 3 {
		t.Errorf("Expected 3 groups, got %d", len(result.Informations))
	}

	if result.Informations[0].Label != "ADCO" {
		t.Errorf("Expected first label to be ADCO, got %s", result.Informations[0].Label)
	}
}

func TestGetSerialConfig(t *testing.T) {
	config := getSerialConfig(ModeHistorical)
	if config.BaudRate != 1200 || config.DataBits != 7 || config.Parity != serial.NoParity || config.StopBits != serial.OneStopBit {
		t.Errorf("Incorrect serial config for historical mode")
	}
}
