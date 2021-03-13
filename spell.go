package spell

// "Development of a Spelling List" by M. Douglas McIlroy
// https://ieeexplore.ieee.org/iel5/26/23952/01095395.pdf
// https://github.com/arnoldrobbins/v10spell

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
