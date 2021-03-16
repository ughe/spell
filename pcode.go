package spell

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type dict struct {
	i    uint16 // (i & 0x07FF) is index of encodes (type []bits)
	word string
}

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

// read an annotated spelling list in form
//	word <tab> affixcode [ , affixcode ] ...
// print a reencoded version
func Pcode() {
	words := make([]dict, 0)
	encodes := make([]bits, 0) // Max size 2^11 (index fits in 11 bits)
	var err error

	if len(os.Args) <= 1 {
		s := bufio.NewScanner(os.Stdin)
		words, encodes, err = readWordEncodings(words, encodes, s)
		if err != nil {
			fatalf("%v\n", err)
		}
	}

	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fatalf("cannot open %s\n%v\n", path, err)
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		words, encodes, err = readWordEncodings(words, encodes, s)
		if err != nil {
			fatalf("%v\n", err)
		}
	}
	fmt.Fprintf(os.Stderr, "words = %d; codes = %d\n", len(words), len(encodes))

	nBytes, err := writeDict(words, encodes, os.Stdout)
	if err != nil {
		fatalf("%v\n", err)
	}
	fmt.Fprintf(os.Stderr, "output bytes = %d\n", nBytes)
}

func readWordEncodings(words []dict, encodes []bits, s *bufio.Scanner) ([]dict, []bits, error) {
	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, nil, fmt.Errorf("Expected 2 words in a line. Found %d for: \"%v\"\n", len(fields), line)
		}
		word := fields[0]
		affixes := fields[1]

		// Find index of bits in encodes, or add if code does not exist
		code, err := strToCode(affixes) // Equivalent of `typecode` and `codetab`
		if err != nil {
			return nil, nil, err
		}
		var i int
		for i = 0; i < len(encodes); i++ {
			if encodes[i] == code {
				break
			}
		}
		if i == len(encodes) {
			encodes = append(encodes, code)
		}

		// Accumulate the encoding index and word
		words = append(words, dict{word: word, i: uint16(i)})
	}
	err := s.Err()

	return words, encodes, err
}

func sread(b *bufio.Reader) (uint16, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(b, buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(buf), nil
}

func lread(b *bufio.Reader) (uint32, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(b, buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf), nil
}

func sput(b *bufio.Writer, bits uint16) error {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, bits)
	_, err := b.Write(buf)
	return err
}

func lput(b *bufio.Writer, bits uint32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, bits)
	_, err := b.Write(buf)
	return err
}

// spit out the encoded dictionary
// all numbers are encoded big-endian.
// struct {
//   ncodes  uint16
//   encodes [ncodes]bits
//   []struct{
//     encode uint16
//     word   []uint16
//   }
// }
// bit mask (for encode uint16) is:
// 0x8000 flag for code word
// 0x7800 count of number of common bytes with previous word
// 0x07ff index into codes array for affixes
func writeDict(words []dict, encodes []bits, w io.Writer) (int, error) {
	sort.Slice(words, func(i, j int) bool {
		return words[i].word < words[j].word
	})

	f := bufio.NewWriter(w)
	defer f.Flush()

	if err := sput(f, uint16(len(encodes))); err != nil {
		return 0, err
	}
	for _, c := range encodes {
		if err := lput(f, uint32(c)); err != nil {
			return 0, err
		}
	}

	nBytes := 2 + 4*len(encodes)
	last := ""
	for _, word := range words {
		var j int
		for j = 0; j < len(word.word) && j < len(last) && word.word[j] == last[j]; j++ {
		}
		if word.word == last {
			fmt.Fprintf(os.Stderr, "identical words: %s\n", word.word)
		}

		// j must fit inside 4 bits. 2^4-1 == 15
		if j > 15 {
			j = 15
		}

		// LSB: Code Index (11 bits) | Common char count (4 bits) | High (1 bit)
		c := (word.i & uint16(0x07FF)) | uint16(((j<<11)&0x7800)|((1<<15)&0x8000))
		if err := sput(f, c); err != nil {
			return nBytes, err
		}

		nBytes += 2

		// Write unique part of word
		if _, err := f.Write([]byte(word.word[j:])); err != nil {
			return nBytes, err
		}

		nBytes += len(word.word[j:])

		last = word.word
	}
	return nBytes, nil
}

// layout of file entry: first byte has bit 0x80 turned on.
// next 4 bits count number of characters common between this
// entry and previous one.  last three bits concatenated with
// second byte are the affixing code, so arranged that the 0x80
// bit is zero in all bytes but the first. 3rd and following
// bytes are the remainder of the dictionary word.
//
// layout in memory: common prefixes are expanded, and the
// first two letters of each word are deleted and found
// instead by lookup in table spacep, which points to the
// first word for each two-letter prefix.
func readDict(path string) ([]dict, [128 * 128]*dict, error) {
	f, err := os.Open(path)
	if err != nil {
		fatalf("spell: cannot open %s\n", path)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	nencode16, err := sread(r)
	if err != nil {
		fatalf("spell: trouble reading %s\n", path)
	}
	nencode := int(nencode16)
	encodes := make([]bits, nencode)
	for i := 0; i < nencode; i++ {
		code, err := lread(r)
		encodes[i] = bits(code)
		if err != nil {
			fatalf("spell: trouble reading %s\n", path)
		}
	}

	space := make([]dict, 0, nencode) // sorted words (first two letters deleted)
	var spacep [128 * 128]*dict       // pointer to words starting with "xx"

	var lasts *string // previous word
	var s *string     // current word

	sp := -1 // previous index into spacep
	i := 0   // current index into spacep

	for err == nil {
		head, err := sread(r) // LSB 11b | 4b | 1b MSB for index and repeated chars
		if err != nil {
			break
		}
		j := (int(head) & 0x7800) >> 11 // num repeated chars (called p in v10 src)

		space = append(space, dict{i: head, word: ""})
		d := &space[len(space)-1]
		s = &d.word

		// Need to update the spacep prefix mapping if the number of repeated
		// characters is less than two. TODO: What happens if word has only 1 char?
		if j <= 0 { // invariant: j >= 0
			c, err := r.ReadByte()
			if err != nil {
				break
			}
			i = int(c) * 128 // Points first x in "xx"
		}
		if j <= 1 {
			c, err := r.ReadByte()
			if err != nil {
				break
			}
			if c&0x80 == 0 { // works since c is ASCII otherwise is 11b|4b|1b enc
				i = i/128*128 + int(c) // Points second x in "xx"
			}
			if i <= sp {
				fatalf("spell: the dict isn't sorted")
			}
			for sp < i {
				sp++
				spacep[sp] = d
			}
		} // else j >= 2 therefore sp already points to correct "xx" location

		// copy repeated chars (can skip first 2 which are "xx" index in spacep)
		*s += (*lasts)[:j-2]

		// copy non-repeated chars
		for {
			c, err := r.ReadByte()
			if err != nil {
				break // breaks two loops since err != nil
			}
			if c&0x80 != 0 { // Not ASCII; part of the encoding
				err = r.UnreadByte()
				break
			}
			*s += string(c)
		}

		lasts = s
	}
	if err != io.EOF {
		fatalf("spell: trouble reading %s\n%v", path, err)
	}

	return space, spacep, nil
}
