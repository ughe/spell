package spell

// "Development of a Spelling List" by M. Douglas McIlroy
// https://ieeexplore.ieee.org/iel5/26/23952/01095395.pdf
// https://github.com/arnoldrobbins/v10spell

import (
	"bufio"
	"flag"
	"os"
	"os/user"
	"path"
	"strings"
)

type prefix struct {
	s    string
	flag bits
}

type suffix struct {
	s         string
	p1        op
	n1        int
	d1        string
	a1        string
	flag      bits
	affixable bits
	p2        op
	n2        int
	d2        string
	a2        string
}

func pair(a, b uint8) bits {
	return (bits(a) << 8) | bits(b)
}

const dLEV = 2
const dSIZ = 40 // Deriv Size

var vflag bool
var xflag bool

// main function for spell (equivalent to main in sprog.c in v10spell)
func Spell() {
	usr, err := user.Current()
	homeDir := "/"
	if err == nil {
		homeDir = usr.HomeDir
	}
	libDir := path.Join(homeDir, "usr", "local", "lib", "v10spell")
	defaultDictPath := path.Join(libDir, "amspell")

	f := flag.String("f", defaultDictPath, "Path to encoded spell dictionary file (created with pcode)")
	v := flag.Bool("v", false, "Print all words not literally in the spelling list, with derivations")
	x := flag.Bool("x", false, "Print on standard error, marked with =, every stem as it is looked up in the spelling list, along with its affix classes. Typically used for maintenance.")
	// Skipping these flags:
	// -b British spelling (can be achieved by using -f brspell)
	// -c Input is one word per line. Outputs + if word known and - if word rejected.
	// -C Input is one word per line. Outputs 0 if word known. Larger numbers indicate words derived by increasingly elaborate paths. Typically used by other programs piping queries to v10spell.

	// Global flags
	vflag = *v
	xflag = *x

	space, spacep := readDict(*f)

	// reset affix

	if len(os.Args) <= 1 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			line := s.Text()
			fields := strings.Fields(line)
			for _, word := range fields {
				// Check each word
			}
		}
		err = s.Err()
		if err != nil {
			fatalf("%v\n", err)
		}
	}

	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fatalf("cannot open %s\n", path)
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			line := s.Text()
			fields := strings.Fields(line)
			// TODO: Check each word
		}
		err = s.Err()
		if err != nil {
			fatalf("%v\n", err)
		}
	}

}
