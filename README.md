# TeleInfo Library

Une bibliothèque Go pour lire les trames TIC (Télé-information Client) des compteurs Linky.
- Lecture des trames TIC via un port série
- Utilisation de channels pour un traitement asynchrone
- Conversion en JSON
- Support mode historique et standard https://www.enedis.fr/media/2035/download

Testé avec le dongle https://github.com/hallard/uTeleinfo

## Installation

```bash
go get github.com/dmachard/go-ticreader
```

## Utilisation

### Example basique

```go
ticChan, err := ticreader.StartReading("/dev/ttyACM0", ticreader.ModeHistorical)
if err != nil {
    fmt.Println("Erreur:", err)
    return
}

for tic := range ticChan {
    teleinfo, _ := tic.ToJSON()
    fmt.Println(teleinfo)
}
```

Example d'une trame au format JSON

```json
{
  "timestamp": "2025-02-19T21:09:37.123405268+01:00",
  "dataset": [
    {
      "label": "ADCO",
      "data": "xxxxxxxxx",
      "valid": true
    },
    {
      "label": "OPTARIF",
      "data": "BASE",
      "valid": true
    },
    {
      "label": "ISOUSC",
      "data": "45",
      "valid": true
    },
    {
      "label": "BASE",
      "data": "xxxxxx",
      "valid": true
    },
    {
      "label": "PTEC",
      "data": "TH..",
      "valid": true
    },
    {
      "label": "IINST",
      "data": "002",
      "valid": true
    },
    {
      "label": "IMAX",
      "data": "090",
      "valid": true
    },
    {
      "label": "PAPP",
      "data": "00530",
      "valid": true
    },
    {
      "label": "HHPHC",
      "data": "A",
      "valid": true
    },
    {
      "label": "MOTDETAT",
      "data": "000000",
      "valid": true
    }
  ]
}
```

## Tests

```bash
go test -v .
```
