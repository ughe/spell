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
	word string
	i    uint16 // index of bits encoding in encodes slice
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
func readDict(path string) ([]dict, []bits, error) {
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

	words := make([]dict, 0, nencode) // At least nencode words
	var spacep [128 * 128]*byte       // pointer to words starting with "xx"
	var space [300000]byte

	var s *byte = &space[0]
	var lasts *byte
	sp := -1 // Index into spacep: Previous
	i := 0   // Index into spacep: Next

	c, err := r.ReadByte()
	for err == nil {
		j := (c >> 3) & 0xF // j aka p aka # repeated characters
		// copy 2 bytes into space buffer
		*s = c
		s++
		if *s, err = r.ReadByte(); err != nil {
			break
		}
		s++

		// Lookup the correct index i in spacep using j (# repeated chars)
		if j == 0 { // invariant: j >= 0
			if c, err = r.ReadByte(); err != nil {
				break
			}
			i = int(c) * 128 // Points first x in "xx" to char c
		} else if j == 1 {
			if c, err = r.ReadByte(); err != nil {
				break
			}
			if c&0x80 == 0 {
				i = i/128*128 + int(c) // Points second x in "xx" to char c
			}
			if i <= sp {
				fatalf("spell: the dict isn't sorted")
			}
			for sp < i {
				sp++
				spacep[sp] = s - 2 // entry starts before last two bytes coppied in
			}
		} // else j >= 2 therefore i already points to correct "xx" location

		ls := lasts // beginning of the last word
		lasts = s   // beginning of the new last word in the space buffer

		// copy repeated chars (can skip first 2 which are "xx" index in spacep)
		for j = j - 2; j > 0; j-- {
			*s = *ls
			s++
			ls++
		}

		// copy non-repeated chars
		for {
			if c, err = r.ReadByte(); err != nil {
				break // breaks two loops since err != nil
			}
			if c & 0x80 {
				break
			}
			*s = c
			s++
		}

		*s = 0
	}
	if err != io.EOF {
		fatalf("spell: trouble reading %s\n%v", path, err)
	}

	for ; sp < len(spacep); sp++ {
		spacep[sp] = s
	}

	return spacep, space, nil
}
