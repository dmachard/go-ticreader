package ticreader

import (
	"testing"

	"go.bug.st/serial"
)

func TestDecodeFrame(t *testing.T) {
	frame := "ADCO 061764523690 H\nOPTARIF BASE 0\nISOUSC 45 ?"
	decodedFrame := decodeFrame(frame, ModeHistorical)

	if len(decodedFrame.Informations) != 3 {
		t.Errorf("Expected 3 groups, got %v", len(decodedFrame.Informations))
	}

	if decodedFrame.Informations[0].Label != "ADCO" {
		t.Errorf("Expected first label to be ADCO, got %s", decodedFrame.Informations[0].Label)
	}
}

func TestGetSerialConfig(t *testing.T) {
	config := getSerialConfig(ModeHistorical)
	if config.BaudRate != 1200 || config.DataBits != 7 || config.Parity != serial.NoParity || config.StopBits != serial.OneStopBit {
		t.Errorf("Incorrect serial config for historical mode")
	}
}
