package ticreader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type GroupInfo struct {
	Label    string `json:"label"`
	Data     string `json:"data"`
	Horodate string `json:"horodate,omitempty"`
	Checksum string `json:"-"`
	Valid    bool   `json:"valid"`
}

type TeleInfo struct {
	Timestamp    time.Time   `json:"timestamp"`
	Informations []GroupInfo `json:"teleinfo"`
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
		for {
			frame, err := readFrame(reader, mode)
			if err != nil {
				log.Error("Error in frame: ", err)
				continue
			}
			frameChan <- frame
		}
	}()

	return frameChan, nil
}

func readFrame(reader *bufio.Reader, mode LinkyMode) (TeleInfo, error) {
	var frameBuilder strings.Builder
	inFrame := false

	for {
		c, err := reader.ReadByte()
		if err != nil {
			return TeleInfo{}, err
		}

		if c == 0x02 {
			frameBuilder.Reset()
			inFrame = true
			continue
		}

		if inFrame {
			if c == 0x03 {
				decodedFrame, err := decodeFrame(frameBuilder.String(), mode)
				if err != nil {
					return TeleInfo{}, err
				}
				return decodedFrame, nil
			}
			frameBuilder.WriteByte(c)
		}
	}
}

func decodeFrame(frame string, mode LinkyMode) (TeleInfo, error) {
	var teleinfo TeleInfo
	lines := strings.Split(frame, "\n")

	for _, line := range lines {
		var group GroupInfo
		var err error
		fmt.Println("decode", line)
		if mode == ModeStandard {
			group, err = parseStandardFrame(line)
		} else {
			group, err = parseHistoricFrame(line)
		}
		fmt.Println(group, err)
		if err != nil {
			return TeleInfo{}, err
		}
		if group.Label != "" {
			teleinfo.Informations = append(teleinfo.Informations, group)
		}
	}
	teleinfo.Timestamp = time.Now()
	return teleinfo, nil
}
