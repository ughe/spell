package spell

// "Development of a Spelling List" by M. Douglas McIlroy
// https://ieeexplore.ieee.org/iel5/26/23952/01095395.pdf
// https://github.com/arnoldrobbins/v10spell

import (
	"encoding/csv"
	"io/ioutil"
)

type bits uint32

func pair(a, b uint8) bits {
	return (bits(a) << 8) | bits(b)
}

// Unknown
const dLEV = 2
const dSIZ = 40

// Prefix
type pretab struct {
	s string
	f int
}

// Suffix
type suftab struct {
	s string
}

func init() {
	// 1. Read prefixes.go CSV
	buf, err := ioutil.ReadFile(imageFilename)
	if err != nil {
		return err
	}
	// 2. Turn them into preftab
	prefixTable := make([][]pretab, int('a')-int('z')+1)
}
