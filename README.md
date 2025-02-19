# TeleInfo Library

Une bibliothèque Go pour lire les trames TIC (Télé-information Client) des compteurs Linky.
- Lecture des trames TIC via un port série
- Utilisation de channels pour un traitement asynchrone
- Conversion en JSON
- Support mode historique et standard

Testé avec le dongle https://github.com/hallard/uTeleinfo

## Installation

```bash
go get github.com/dmachard/go-teleinfolib
```

## Utilisation

### Example basique

```go
frameChan, err := teleinfolib.StartReading("/dev/ttyACM0", teleinfolib.ModeHistorical)
if err != nil {
    fmt.Println("Erreur:", err)
    return
}

for frame := range frameChan {
    fmt.Println(frame.ToJSON())
}
```

## Tests

```bash
go test -v .
```
