package spell

import (
	"bytes"
	_ "embed"
	"encoding/csv"
)

//go:embed prefix.csv
var prefixCSV []byte

//go:embed suffix.csv
var suffixCSV []byte

func init() {
	// 1. Read prefixes.go CSV
	r := csv.NewReader(bytes.NewBuffer(prefixCSV))

	// 2. Turn them into preftab
	prefixTable := make([][]pretab, int('a')-int('z')+1)
}
