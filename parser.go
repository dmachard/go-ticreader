package ticreader

import (
	"errors"
	"strings"
)

var ErrInvalidFrame = errors.New("invalid frame")

// parseHistoricFrame extrait l'étiquette, la donnée et le checksum d'un groupe en mode historique
func parseHistoricFrame(frame string) (GroupInfo, error) {
	parts := strings.Fields(frame)
	if len(parts) < 3 {
		return GroupInfo{}, ErrInvalidFrame
	}

	label := parts[0]
	data := parts[1]
	checksum := parts[2]

	valid := verifyChecksum(label, data, checksum)
	return GroupInfo{Label: label, Data: data, Valid: valid}, nil
}

// parseStandardFrame extrait les données d'une ligne en mode STANDARD
func parseStandardFrame(frame string) (GroupInfo, error) {
	parts := strings.Split(frame, "\t") // Les champs sont séparés par HT (0x09)

	if len(parts) == 4 { // Avec horodatage
		label := parts[0]
		horodate := parts[1]
		data := parts[2]
		checksum := parts[3]

		valid := verifyChecksum(label+horodate, data, checksum)
		return GroupInfo{Label: label, Horodate: horodate, Data: data, Valid: valid}, nil

	} else if len(parts) == 3 { // Sans horodatage
		label := parts[0]
		data := parts[1]
		checksum := parts[2]

		valid := verifyChecksum(label, data, checksum)
		return GroupInfo{Label: label, Data: data, Valid: valid}, nil
	}

	return GroupInfo{}, ErrInvalidFrame
}

// verifyChecksum vérifie la validité du checksum
func verifyChecksum(label, value, checksum string) bool {
	expectedChecksum := calculateChecksum(label, value)
	return checksum == string(expectedChecksum)
}

// calculateChecksum calcule le checksum attendu
func calculateChecksum(label, value string) byte {
	chksum := 32

	for _, c := range label {
		chksum += int(c)
	}
	for _, c := range value {
		chksum += int(c)
	}

	chksum = (chksum & 63) + 32
	return byte(chksum)
}
