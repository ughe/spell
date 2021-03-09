package spell

// affix codes

const ED	= (1<<0) /* +ed, +ing */
const ADJ	= (1<<1) /* (nce)-t_ce, +ize,+al, +ness, -t+cy, +ity, +ly */

const NOUN	= (1<<2) /* +s (+es), +make, +hood, +ship +less */
const PROP_COLLECT = (1<<3) /* +'s, +an, +ship(for -manship) +less */
const ACTOR	= (1<<4) /* +er */
const EST	= (1<<5)
const COMP	= (EST|ACTOR) /* +er,+est */
const DONT_TOUCH = (1<<6)
const ION	= (1<<7) /* +ion, +or */
const N_AFFIX	= (1<<8) /* +ic, +ive, +ize, +like, +al, +ful, +ism, +ist, -t+cy, +c (maniac) */
const V_AFFIX	= (1<<9) /* +able, +ive, +ity((bility), +ment */
const V_IRREG	= (1<<10) /* +ing +es +s*/
const VERB	= (V_IRREG|ED)
const MAN	= (1<<11) /* +man, +men, +women, +woman */
const ADV	= (1<<12) /* +hood, +ness */
const STOP	= (1<<14) /* stop list */

const NOPREF	= (1<<13) /* no prefix */
const MONO	= (1<<15) /* double final consonant as in fib->fibbing */
const IN	= (1<<16) /* in- im- ir, not un- */
const _Y	= (1<<17) /* +y */

const ALL	= (~(NOPREF|STOP|DONT_TOUCH|MONO|IN)) /*anything goes (no stop or nopref)*/
