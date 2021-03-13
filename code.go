package spell

import (
	"fmt"
	"strings"
)

// affix codes

type (
	bits uint32
	op   func(ep, dd, a string, lev, flag uint32) bits
)

const (
	ED           bits = 1 << 0 /* +ed, +ing */
	ADJ          bits = 1 << 1 /* (nce)-t_ce, +ize,+al, +ness, -t+cy, +ity, +ly */
	NOUN         bits = 1 << 2 /* +s (+es), +make, +hood, +ship +less */
	PROP_COLLECT bits = 1 << 3 /* +'s, +an, +ship(for -manship) +less */
	ACTOR        bits = 1 << 4 /* +er */
	EST          bits = 1 << 5
	COMP         bits = EST | ACTOR /* +er,+est */
	DONT_TOUCH   bits = 1 << 6
	ION          bits = 1 << 7  /* +ion, +or */
	N_AFFIX      bits = 1 << 8  /* +ic, +ive, +ize, +like, +al, +ful, +ism, +ist, -t+cy, +c (maniac) */
	V_AFFIX      bits = 1 << 9  /* +able, +ive, +ity(bility), +ment */
	V_IRREG      bits = 1 << 10 /* +ing +es +s*/
	VERB         bits = V_IRREG | ED
	MAN          bits = 1 << 11 /* +man, +men, +women, +woman */
	ADV          bits = 1 << 12 /* +hood, +ness */
	NOPREF       bits = 1 << 13 /* no prefix */
	STOP         bits = 1 << 14 /* stop list */
	MONO         bits = 1 << 15 /* double final consonant as in fib->fibbing */
	IN           bits = 1 << 16 /* in- im- ir, not un- */
	Y            bits = 1 << 17 /* +y */ // Formerly _Y

	ALL bits = ^(NOPREF | STOP | DONT_TOUCH | MONO | IN) /*anything goes (no stop or nopref)*/
)

// bits identifiers used in initial dictionaries. Similar to `codetab` in c version
var nameCodes map[string]bits = map[string]bits{
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
	"y":      Y,
}

// Converts a comma-separated string of bits identifiers to bits
func strToCode(s string) (bits, error) {
	var code bits = 0
	for _, name := range strings.Split(strings.TrimSpace(s), ",") {
		b, ok := nameCodes[name]
		if !ok {
			return code, fmt.Errorf("Unknown affix code \"%s\"\n", name)
		}
		code = code | b
	}
	return code, nil
}

// Outputs | separated string of bits names (excluding composite ones)
func codeToStr(code bits) string {
	buf := ""
	for k, v := range codeNames {
		if code&k != 0 {
			buf += "|" + v
		}
	}
	return buf[1:] // Strip leading '|'
}

// Excludes composit COMP, ALL, and VERB names
var codeNames map[bits]string = map[bits]string{
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

var opCodes map[string]op = map[string]op{
	"nop":    nop,
	"strip":  strip,
	"ize":    ize,
	"i_to_y": i_to_y,
	"ily":    ily,
	"subst":  subst,
	"CCe":    CCe,
	"tion":   tion,
	"an":     an,
	"s":      s,
	"es":     es,
	"bility": bility,
	"y_to_e": y_to_e,
	"VCe":    VCe,
}

// Returns true if character is a vowel
// Implementation notes: Original C uses a LUT of [128]bool
// Instead of for loop, could also use map, switch, boolean or
func isVowel(c rune) bool {
	vowels := "aeiouyAEIOUY"
	for _, s := range vowels {
		if c == s {
			return true
		}
	}
	return false
}

func nop(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func cstrip(ep, dd, a string, lev, flag uint32) bits {
	/*
		if isVowel(ep[0]) && isVowel(ep[-1]) {
			// TODO:
			// if (ep[-1], ep[0]) ==
			// a,a
			// a,e
			// a,i
			// e,a
			// e,e
			// e,i
			// i,i
			// o,a
			return 0
		} else if ep[0] == ep[-1] && ep[0] == ep[-2] {
			return 0
		} else {
			return strip(ep, d, a, lev, flag)
		}
	*/
	return 0
}

func isSet(a, b bits) bool {
	return (a & b) != 0
}

func strip(ep, dd, a string, lev, flag uint32) bits {
	/*
		var h bits = trypref(ep, a, lev, flag)
		if !(isSet(h, MONO) && isVowel(ep[0]) && isVowel(ep[-2])) {
			return h
		}
		if isVowel(ep[0]) && !isVowel(ep[-1]) && ep[-1] == ep[-2] {
			h = trypref(ep-1, a, lev, flag|MONO)
			if h != 0 {
				return h
			}
			return trysuff(ep, lev, flag)
		}
	*/
	return 0
}

func ize(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func i_to_y(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func ily(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func subst(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func CCe(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func tion(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func an(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func s(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func es(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func bility(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func y_to_e(ep, dd, a string, lev, flag uint32) bits {
	return 0
}

func VCe(ep, dd, a string, lev, flag uint32) bits {
	return 0

}
