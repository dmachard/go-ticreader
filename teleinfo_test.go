package ticreader

import (
	"testing"

	"go.bug.st/serial"
)

func TestGetSerialConfig(t *testing.T) {
	config := getSerialConfig(ModeHistorical)
	if config.BaudRate != 1200 || config.DataBits != 7 || config.Parity != serial.NoParity || config.StopBits != serial.OneStopBit {
		t.Errorf("Incorrect serial config for historical mode")
	}
}
