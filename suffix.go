package spell

// Return func to order reversed suffixes alphabetically
// For example: sort.Slice(suffixes, orderSuffix(suffixes))

// Returns less func for sort.Slice Suffixes are ordered
// by reversing each string and comparing them alphabetically
// i.e. (phobia, ac, istic, itic) => (aibohp, ca, citsi, citi)
func orderSuffix(suf []suffix) func(i, j int) bool {
	return func(i, j int) bool {
		a, b := suf[i].s, suf[j].s
		l := len(a) // min(len(a), len(b))
		if len(b) < l {
			l = len(b)
		}
		// Order suffixes backwards (reversed string)
		for k := 1; k <= l; k++ {
			if a[len(a)-k] != b[len(b)-k] {
				return a[len(a)-k] < b[len(b)-k]
			}
		}
		// if suf[i].s[-k:] == suf[j].s[-k:] ; Return longer
		return len(a) > len(b)
	}
}

var suffixes []suffix = []suffix{
	{"phobia", subst, 1, "-e+ia", "", NOUN, NOUN,
		nop, 0, "", ""},

	{"ac", strip, 1, "", "+c", N_AFFIX, ADJ | NOUN,
		nop, 0, "", ""},

	{"istic", strip, 2, "", "+ic", N_AFFIX, ADJ | N_AFFIX | NOUN,
		nop, 0, "", ""},

	{"itic", ize, 1, "-e+ic", "", N_AFFIX, ADJ,
		nop, 0, "", ""},

	{"graphic", i_to_y, 1, "-y+ic", "", NOUN, ADJ | NOUN,
		nop, 0, "", ""},

	{"scopic", ize, 1, "-e+ic", "", NOUN, ADJ,
		nop, 0, "", ""},

	{"metric", i_to_y, 1, "-y+ic", "", NOUN, ADJ,
		nop, 0, "", ""},

	{"logic", i_to_y, 1, "-y+ic", "", NOUN, ADJ,
		nop, 0, "", ""},

	{"onomic", i_to_y, 1, "-y+ic", "", NOUN, ADJ,
		nop, 0, "", ""},

	{"phobic", subst, 1, "-e+ic", "", NOUN, ADJ,
		nop, 0, "", ""},

	{"ed", strip, 1, "", "+d", ED, ADJ | COMP, i_to_y, 2, "-y+ied", "+ed"},

	{"hood", ily, 4, "-y+ihood", "+hood", NOUN | ADV, NOUN,
		nop, 0, "", ""},

	{"nce", subst, 1, "-t+ce", "", ADJ, N_AFFIX | Y | NOUN | VERB | ACTOR | V_AFFIX,
		nop, 0, "", ""},

	{"faible", i_to_y, 4, "-y+iable", "", V_IRREG, ADJ,
		nop, 0, "", ""},

	{"able", CCe, 4, "-e+able", "+able", V_AFFIX, ADJ,
		nop, 0, "", ""},

	{"ive", subst, 0, "-ion+ive", "", N_AFFIX | V_AFFIX, NOUN | N_AFFIX | ADJ,
		nop, 0, "", ""},

	{"ize", CCe, 3, "-e+ize", "+ize", N_AFFIX | ADJ, V_AFFIX | VERB | ION | COMP,
		nop, 0, "", ""},

	{"like", strip, 4, "", "+like", N_AFFIX, ADJ,
		nop, 0, "", ""},

	{"eeing", strip, 3, "", "+ing", V_IRREG, ADJ | NOUN,
		nop, 0, "", ""},

	{"making", strip, 6, "", "+making", NOUN, NOUN,
		nop, 0, "", ""},

	{"keeping", strip, 7, "", "+keeping", NOUN, NOUN,
		nop, 0, "", ""},

	{"ing", CCe, 3, "-e+ing", "+ing", V_IRREG, ADJ | ED | NOUN,
		nop, 0, "", ""},

	{"oidal", strip, 2, "", "+al", NOUN | ADJ, ADJ,
		nop, 0, "", ""},

	{"ical", strip, 2, "", "+al", NOUN | ADJ, ADJ | NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"mental", strip, 2, "", "+al", N_AFFIX, ADJ,
		nop, 0, "", ""},

	{"ional", strip, 2, "", "+al", N_AFFIX, ADJ | NOUN,
		nop, 0, "", ""},

	{"ful", ily, 3, "-y+iful", "+ful", N_AFFIX, ADJ | NOUN,
		nop, 0, "", ""},

	{"ism", CCe, 3, "-e+ism", "ism", N_AFFIX | ADJ, NOUN,
		nop, 0, "", ""},

	{"ogram", subst, -1, "-ph+m", "", NOUN, NOUN,
		nop, 0, "", ""},

	{"ification", i_to_y, 6, "-y+ication", "", ION, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"ization", ize, 4, "-e+ation", "", ION, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"tion", tion, 3, "-e+ion", "+ion", ION, NOUN | N_AFFIX | V_AFFIX | VERB | ACTOR,
		nop, 0, "", ""},

	{"onian", an, 3, "", "+ian", NOUN | PROP_COLLECT, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"woman", strip, 5, "", "+woman", MAN, PROP_COLLECT | N_AFFIX,
		nop, 0, "", ""},

	{"man", strip, 3, "", "+man", MAN, PROP_COLLECT | N_AFFIX | VERB,
		nop, 0, "", ""},

	{"an", an, 1, "", "+n", NOUN | PROP_COLLECT, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"women", strip, 5, "", "+women", MAN, PROP_COLLECT,
		nop, 0, "", ""},

	{"men", strip, 3, "", "+man", MAN, PROP_COLLECT,
		nop, 0, "", ""},

	{"ship", strip, 4, "", "+ship", NOUN | PROP_COLLECT, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"grapher", subst, 1, "-y+er", "", ACTOR, NOUN, strip, 2, "", "+er"},

	{"graphyer", nop, 0, "", "", 0, NOUN,
		nop, 0, "", ""},

	{"maker", strip, 5, "", "+maker", NOUN, NOUN,
		nop, 0, "", ""},

	{"keeper", strip, 6, "", "+keeper", NOUN, NOUN,
		nop, 0, "", ""},

	{"er", strip, 1, "", "+r", ACTOR, NOUN | N_AFFIX | VERB | ADJ, i_to_y, 2, "-y+ier", "+er"},

	{"ator", tion, 2, "-e+or", "", ION, NOUN | N_AFFIX | Y,
		nop, 0, "", ""},

	{"ctor", tion, 2, "", "+or", ION, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"ptor", tion, 2, "", "+or", ION, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"ness", ily, 4, "-y+iness", "+ness", ADJ | ADV, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"less", ily, 4, "-y+iless", "+less", NOUN | PROP_COLLECT, ADJ,
		nop, 0, "", ""},

	{"es", s, 1, "", "+s", NOUN | V_IRREG, DONT_TOUCH, es, 2, "-y+ies", "+es"},

	{"'s", s, 2, "", "+'s", PROP_COLLECT | NOUN, DONT_TOUCH,
		nop, 0, "", ""},

	{"s", s, 1, "", "+s", NOUN | V_IRREG, DONT_TOUCH,
		nop, 0, "", ""},

	{"ment", strip, 4, "", "+ment", V_AFFIX, NOUN | N_AFFIX | ADJ | VERB,
		nop, 0, "", ""},

	{"est", strip, 2, "", "+st", EST, DONT_TOUCH, i_to_y, 3, "-y+iest", "+est"},

	{"logist", i_to_y, 2, "-y+ist", "", N_AFFIX, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"ist", CCe, 3, "-e+ist", "+ist", N_AFFIX | ADJ, NOUN | N_AFFIX | COMP,
		nop, 0, "", ""},

	{"blity", nop, 0, "", "", 0, NOUN,
		nop, 0, "", ""},

	{"ncy", subst, 1, "-t+cy", "", ADJ | N_AFFIX, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"bility", bility, 5, "-le+ility", "", ADJ | V_AFFIX, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"ousity",
		nop, 0, "", "", NOUN, 0, nop, 0, "", ""},

	{"ity", CCe, 3, "-e+ity", "+ity", ADJ, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"bly", y_to_e, 1, "-e+y", "", ADJ, ADV,
		nop, 0, "", ""},

	{"cly", nop, 0, "", "", 0, 0,
		nop, 0, "", ""},

	{"ly", ily, 2, "-y+ily", "+ly", ADJ, ADV | COMP,
		nop, 0, "", ""},

	{"metry", subst, 0, "-er+ry", "", NOUN, NOUN | N_AFFIX,
		nop, 0, "", ""},

	{"y", CCe, 1, "-e+y", "+y", Y, ADJ | COMP,
		nop, 0, "", ""},
}
