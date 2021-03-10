package spell

import (
	"fmt"
	"strings"
)

// affix codes

type Bits uint32

const (
	ED           Bits = 1 << 0 /* +ed, +ing */
	ADJ          Bits = 1 << 1 /* (nce)-t_ce, +ize,+al, +ness, -t+cy, +ity, +ly */
	NOUN         Bits = 1 << 2 /* +s (+es), +make, +hood, +ship +less */
	PROP_COLLECT Bits = 1 << 3 /* +'s, +an, +ship(for -manship) +less */
	ACTOR        Bits = 1 << 4 /* +er */
	EST          Bits = 1 << 5
	COMP         Bits = EST | ACTOR /* +er,+est */
	DONT_TOUCH   Bits = 1 << 6
	ION          Bits = 1 << 7  /* +ion, +or */
	N_AFFIX      Bits = 1 << 8  /* +ic, +ive, +ize, +like, +al, +ful, +ism, +ist, -t+cy, +c (maniac) */
	V_AFFIX      Bits = 1 << 9  /* +able, +ive, +ity(bility), +ment */
	V_IRREG      Bits = 1 << 10 /* +ing +es +s*/
	VERB         Bits = V_IRREG | ED
	MAN          Bits = 1 << 11 /* +man, +men, +women, +woman */
	ADV          Bits = 1 << 12 /* +hood, +ness */
	NOPREF       Bits = 1 << 13 /* no prefix */
	STOP         Bits = 1 << 14 /* stop list */
	MONO         Bits = 1 << 15 /* double final consonant as in fib->fibbing */
	IN           Bits = 1 << 16 /* in- im- ir, not un- */
	Y            Bits = 1 << 17 /* +y */ // Formerly _Y (Go needs Y instead)

	ALL Bits = ^(NOPREF | STOP | DONT_TOUCH | MONO | IN) /*anything goes (no stop or nopref)*/
)

// Bits identifiers used in initial dictionaries
var nameCodes map[string]Bits = map[string]Bits{
	"a":      ADJ,
	"adv":    ADV,
	"comp":   COMP,
	"d":      DONT_TOUCH,
	"ed":     ED,
	"er":     ACTOR,
	"in":     IN,
	"ion":    ION,
	"man":    MAN,
	"ms":     MONO,
	"n":      NOUN,
	"na":     N_AFFIX,
	"nopref": NOPREF,
	"pc":     PROP_COLLECT,
	"s":      STOP,
	"v":      VERB,
	"va":     V_AFFIX,
	"vi":     V_IRREG,
	"y":      _Y,
}

// Converts a comma-separated string of Bits identifiers to Bits
func strToCode(s string) (Bits, error) {
	code := Bits{0}
	for _, name := range strings.Split(strings.Trim(s), ",") {
		b, err := nameCodes[name]
		if err != nil {
			return code, fmt.Errorf("Unknown affix code \"%s\"\n%v", name, err)
		}
		code = code | b
	}
	return code, nil
}

// Outputs | separated string of Bits names (excluding composite ones)
func codeToStr(code Bits) string {
	buf := ""
	for k, v := range codeNames {
		if code&k != 0 {
			buf += "|" + v
		}
	}
	return buf[1:] // Strip leading '|'
}

// Excludes composit COMP, ALL, and VERB names
var codeNames map[Bits]string = map[Bits]string{
	ED:           "ED",
	ADJ:          "AD",
	NOUN:         "NOUN",
	PROP_COLLECT: "PROP_COLLECT",
	ACTOR:        "ACTOR",
	EST:          "EST",
	DONT_TOUCH:   "DONT_TOUCH",
	ION:          "ION",
	N_AFFIX:      "N_AFFIX",
	V_AFFIX:      "V_AFFIX",
	V_IRREG:      "V_IRREG",
	MAN:          "MAN",
	ADV:          "ADV",
	STOP:         "STOP",
	NOPREF:       "NOPREF",
	MONO:         "MONO",
	IN:           "IN",
	Y:            "_Y", // Formerly _Y (Go needs Y instead)
}
