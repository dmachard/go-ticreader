package ticreader

import (
	"bufio"
	"errors"
	"strings"
	"time"
)

var ErrInvalidHistoricalDataset = errors.New("invalid historical dataset")
var ErrInvalidStandardDataset = errors.New("invalid standard dataset")

// Format de la trame
// Une trame est constituée de trois parties
// | STX | Data set | Data set | …. | Data set | ETX
// le caractère "Start TeXt" STX (0x02) indique le début de la trame
// le corps de la trame est composé de plusieurs groupes d'informations,
// le caractère "End TeXt" ETX (0x03) indique la fin de la trame.
func decodeFrame(reader *bufio.Reader, mode LinkyMode) (TeleInfo, error) {
	var frameBuilder strings.Builder
	inFrame := false

	for {
		c, err := reader.ReadByte()
		if err != nil {
			return TeleInfo{ErrorMsg: err.Error()}, err
		}

		if c == 0x02 {
			frameBuilder.Reset()
			inFrame = true
			continue
		}

		if inFrame {
			if c == 0x03 {
				return decodeDataset(frameBuilder.String(), mode), nil
			}
			frameBuilder.WriteByte(c)
		}
	}
}

// Format des groupes d’information
// un caractère "Line Feed" LF (0x0A) indiquant le début du groupe,
// un caractère "Carriage Return" CR (0x0D) indiquant la fin du groupe d'information
func decodeDataset(frame string, mode LinkyMode) TeleInfo {
	var teleinfo TeleInfo
	var currentDataset strings.Builder
	newDataset := false

	for _, c := range frame {
		switch c {
		case 0x0A: // LF: Début d'un nouveau groupe d'information
			currentDataset.Reset()
			newDataset = true

		case 0x0D: // CR: Fin du groupe d'information
			if newDataset {
				var dataset Dataset
				var err error

				if mode == ModeStandard {
					dataset, err = parseStandardDataset(currentDataset.String())
				} else {
					dataset, err = parseHistoricDataset(currentDataset.String())
				}

				if err != nil {
					return TeleInfo{ErrorMsg: err.Error(), ErrorDetails: currentDataset.String()}
				}

				if dataset.Label != "" {
					teleinfo.Dataset = append(teleinfo.Dataset, dataset)
				}
			}
			newDataset = false

		default:
			if newDataset {
				currentDataset.WriteRune(c)
			}
		}
	}

	teleinfo.Timestamp = time.Now()
	return teleinfo
}

// Format
// | Etiquette | Separator | Donnée | Separator | Checksum |
// le champ étiquette dont la longueur est inférieure ou égale à huit caractères,
// le champ "donnée" dont la longueur est variable
// Le séparateur est un espace SP (0x20) en mode historique et une tabulation HT (0x09) en mode standard

// parseHistoricDataset extrait l'étiquette, la donnée et le checksum d'un groupe en mode historique
func parseHistoricDataset(dataset string) (Dataset, error) {
	parts := strings.Fields(dataset)
	if len(parts) < 3 {
		return Dataset{}, ErrInvalidHistoricalDataset
	}

	label := parts[0]
	data := parts[1]
	checksum := parts[2]

	valid := verifyChecksum(label+" "+data, checksum)
	return Dataset{Label: label, Data: data, Valid: valid}, nil
}

// parseStandardFrame extrait les données d'une ligne en mode STANDARD
func parseStandardDataset(dataset string) (Dataset, error) {
	parts := strings.Split(dataset, "\t") // Les champs sont séparés par HT (0x09)

	if len(parts) == 4 { // Avec horodatage
		label := parts[0]
		horodate := parts[1]
		data := parts[2]
		checksum := parts[3]

		valid := verifyChecksum(checksum, label+"\t"+horodate+"\t"+data)
		return Dataset{Label: label, Horodate: horodate, Data: data, Valid: valid}, nil

	} else if len(parts) == 3 { // Sans horodatage
		label := parts[0]
		data := parts[1]
		checksum := parts[2]

		valid := verifyChecksum(checksum, label+"\t"+data)
		return Dataset{Label: label, Data: data, Valid: valid}, nil
	}

	return Dataset{}, ErrInvalidStandardDataset
}

// verifyChecksum vérifie la validité du checksum
func verifyChecksum(checksum string, data string) bool {
	expectedChecksum := calculateChecksum(data)
	return checksum == string(expectedChecksum)
}

// calculateChecksum calcule le checksum attendu
// La checksum est calculée sur l'ensemble des caractères allant du début du champ Etiquette
// à la fin du champ Donnée, séparateurs inclus.
// Le résultat sera toujours un caractère ASCII imprimable compris entre 0x20 et 0x5F
func calculateChecksum(data string) byte {
	sum := 0

	// Additionne tous les caractères du début du champ "Étiquette" jusqu'au séparateur HT (0x09)
	for _, c := range data {
		sum += int(c)
	}

	// Tronque sur 6 bits (AND 0x3F) et ajoute 0x20 pour obtenir un caractère imprimable
	return byte((sum & 0x3F) + 0x20)
}
