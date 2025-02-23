package ticreader

import (
	"bufio"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type Dataset struct {
	Label    string `json:"label"`
	Data     string `json:"data"`
	Horodate string `json:"horodate,omitempty"`
	Checksum string `json:"-"`
	Valid    bool   `json:"valid"`
}

type TeleInfo struct {
	Timestamp    time.Time `json:"timestamp"`
	Dataset      []Dataset `json:"dataset"`
	ErrorMsg     string    `json:"error-msg,omitempty"`
	ErrorDetails string    `json:"error-details,omitempty"`
}

func (t TeleInfo) ToJSON() (string, error) {
	jsonData, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type LinkyMode struct {
	BaudRate int
	DataBits int
	Parity   serial.Parity
	StopBits serial.StopBits
}

var (
	ModeStandard   = LinkyMode{9600, 7, serial.EvenParity, serial.OneStopBit}
	ModeHistorical = LinkyMode{1200, 7, serial.NoParity, serial.OneStopBit}
)

func getSerialConfig(mode LinkyMode) *serial.Mode {
	return &serial.Mode{
		BaudRate: mode.BaudRate,
		DataBits: mode.DataBits,
		Parity:   mode.Parity,
		StopBits: mode.StopBits,
	}
}

func StartReading(port string, mode LinkyMode) (<-chan TeleInfo, error) {
	serialConfig := getSerialConfig(mode)
	stream, err := serial.Open(port, serialConfig)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(stream)
	frameChan := make(chan TeleInfo)

	go func() {
		defer close(frameChan)

		for {
			frame, err := decodeFrame(reader, mode)
			if err != nil {
				log.Error("Error to read serial: ", err)
				return
			}
			frameChan <- frame
		}
	}()

	return frameChan, nil
}
