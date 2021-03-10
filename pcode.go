// read an annotated spelling list in form
//	word <tab> affixcode [ , affixcode ] ...
// print a reencoded version
package spell

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"sort"
	"strings"
)

type dict struct {
	word string
	i    uint16 // index of Bits encoding in encodes slice
}

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a)
	os.Exit(1)
}

func main() {
	words := make([]dict, 0, 200000)
	encodes := make([]Bits, 0, 2048) // Max size 2^11 (11 bit binary format)
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
			fatalf("Cannot open %s\n%v\n", path, err)
		}
		s := bufio.NewScanner(f)
		words, encodes, err = readWordEncodings(words, encodes, s)
		if err != nil {
			fatalf("%v\n", err)
		}
	}

	if err := writeDict(words, encodes); err != nil {
		fatalf("%v\n", err)
	}
}

func readWordEncodings(words []dict, encodes []Bits, s *bufio.Scanner) ([]dict, []Bits, error) {

	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, nil, fmt.Errorf("Expected 2 words in a line. Found %d for: \"%v\"\n", len(fields), line)
		}
		word := fields[0]
		affixes := fields[1]

		// Find index of Bits in encodes, or add if code does not exist
		code, err := strToCode(affixes)
		if err != nil {
			return nil, nil, err
		}
		var c Bits
		var i int
		for i, c = range encodes {
			if c == code {
				break
			}
		}
		if i == len(encodes) {
			encodes = append(encodes, code)
		}

		// Accumulate the word and encoding index
		words = append(words, dict{word: word, i: uint16(i)})
	}
	err := s.Err()

	fmt.Fprintf(os.Stderr, "words = %d; codes = %d\n", len(words), len(encodes))
	return words, encodes, err
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
//	struct
//	{
//		short	ncodes;
//		int	encodes[ncodes];
//		struct
//		{
//			short	encode;
//			char	word[*];
//		} words[*];
//	};
// 0x8000 flag for code word
// 0x7800 count of number of common bytes with previous word
// 0x07ff index into codes array for affixes
func writeDict(words []dict, encodes []Bits) error {
	sort.Slice(words, func(i, j int) bool {
		return words[i].word < words[j].word
	})

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()

	if err := sput(f, uint16(len(encodes))); err != nil {
		return err
	}
	for _, c := range encodes {
		if err := lput(f, uint32(c)); err != nil {
			return err
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

		// Code Index (11 bits) | Common char count (4 bits) | High (1 bit)
		c := (word.i & uint16(0x07FF)) | uint16(((j<<11)&0x7800)|((1<<15)&0x8000))
		if err := sput(f, c); err != nil {
			return nil
		}

		nBytes += 2

		// Write unique part of word
		if _, err := f.Write([]byte(word.word[j:])); err != nil {
			return err
		}

		last = word.word
	}
	return nil
}
