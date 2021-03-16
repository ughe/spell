package spell

import (
	"bytes"
	"testing"
)

func TestWriteDict(t *testing.T) {
	paths := []string{
		"dictionaries/list",
		"dictionaries/american",
		"dictionaries/local",
		"dictionaries/stop",
	}

	words := make([]dict, 0)
	encodes := make([]bits, 0)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("cannot open %s\n", path)
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		words, encodes, err = readWordEncodings(words, encodes, s)
		if err != nil {
			fatalf("%v\n", err)
		}
	}

	var buf bytes.Buffer
	nBytes, err := writeDict(words, encodes, &buf)
	if err != nil {
		t.Fatalf("failed writeDict with err: %v", err)
	}
	if nBytes != len(buf) {
		t.Fatalf("writeDict nBytes != len(buf). (%d != %d)", nBytes, len(buf))
	}

	wordsPrime, encodesPrime, err := readDictTesting(&buf)
	if !reflect.DeepEqual(words, wordsPrime) {
		t.Fatalf("input of writeDict != output of readDictTesting (words)")
	}
	if !reflect.DeepEqual(encodes, encodesPrime) {
		t.Fatalf("input of writeDict != output of readDictTesting (encodes)")
	}
}

// Returns words and encodings given the output of writeDict
// Performs the opposite operation of writeDict. For testing purposes
func readDictTesting(rd io.Reader) ([]dict, []bits, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
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
	var last string = ""

	for {
		c, err := sread(r)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fatalf("spell: trouble reading %s\n", path)
			}
		}
		i := uint16(c & 0x07FF) // encodes index lookup
		j := (c & 0x7800) >> 11 // length of common prefix with previous word
		if c&0x8000 == 0 {
			fatalf("spell: trouble reading %s (invariant violation)\n", path)
		}

		remainder := ""
		b, err := r.ReadByte()
		for ; b&0x8000 == 0 && err == nil; b, err = r.ReadByte() {
			remainder += rune(b)
		}
		if err == nil && err != io.EOF {
			fatalf("spell: trouble reading %s\n", path)
		}
		if b&0x8000 != 0 {
			if err := r.UnreadByte(); err != nil {
				fatalf("spell: trouble reading %s (%v)\n", path, err)
			}
		}

		word := last[:j] + remainder
		words = append(words, dict{word: word, i: i})
		last = word
	}

	return words, encodes, nil
}
