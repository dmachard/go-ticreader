package ticreader

import (
	"bufio"
	"encoding/json"
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
	Timestamp          time.Time   `json:"timestamp"`
	Informations       []GroupInfo `json:"teleinfo"`
	DecodeErrorMsg     string      `json:"decode-error-msg,omitempty"`
	DecodeErrorDetails string      `json:"decode-error-details,omitempty"`
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

// Format de la trame
// Une trame est constituée de trois parties
// | STX | Data set | Data set | …. | Data set | ETX
// le caractère "Start TeXt" STX (0x02) indique le début de la trame
// le corps de la trame est composé de plusieurs groupes d'informations,
// le caractère "End TeXt" ETX (0x03) indique la fin de la trame.
func readFrame(reader *bufio.Reader, mode LinkyMode) (TeleInfo, error) {
	var frameBuilder strings.Builder
	inFrame := false

	for {
		c, err := reader.ReadByte()
		if err != nil {
			return TeleInfo{DecodeErrorMsg: err.Error()}, err
		}

		if c == 0x02 {
			frameBuilder.Reset()
			inFrame = true
			continue
		}

		if inFrame {
			if c == 0x03 {
				return decodeFrame(frameBuilder.String(), mode), nil
			}
			frameBuilder.WriteByte(c)
		}
	}
}

// Format des groupes d’information
// un caractère "Line Feed" LF (0x0A) indiquant le début du groupe,
// un caractère "Carriage Return" CR (0x0D) indiquant la fin du groupe d'information
func decodeFrame(frame string, mode LinkyMode) TeleInfo {
	var teleinfo TeleInfo
	lines := strings.Split(frame, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var group GroupInfo
		var err error
		if mode == ModeStandard {
			group, err = parseStandardFrame(line)
		} else {
			group, err = parseHistoricFrame(line)
		}
		if err != nil {
			return TeleInfo{DecodeErrorMsg: err.Error(), DecodeErrorDetails: line}
		}
		if group.Label != "" {
			teleinfo.Informations = append(teleinfo.Informations, group)
		}
	}
	teleinfo.Timestamp = time.Now()
	return teleinfo
}
