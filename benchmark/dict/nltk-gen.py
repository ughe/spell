#!/usr/bin/env python3
import nltk

def isword(word):
    return all([x >= 'a' and x <= 'z' for x in word.lower()])

def dumpwords(allwords, dst, dstother=None):
    allwords = [str(x) for x in allwords]
    words = sorted(filter(lambda x: x, [x if isword(x) else None for x in allwords]), key=str.casefold)
    other = sorted(filter(lambda x: x, [x if not isword(x) else None for x in allwords]), key=str.casefold)
    with open(dst, 'w') as f:
        f.write('\n'.join(words))
    if dstother is not None:
        with open(dstother, 'w') as f:
            f.write('\n'.join(other))
    print('Created %s\nTotal: %d. Words: %d. Invalid: %d.' % (dst, len(allwords), len(words), len(other)))

# wordnet.txt
nltk.download('wordnet')
dumpwords(nltk.corpus.wordnet.all_lemma_names(), 'wordnet.txt')

# /usr/share/dict/words
nltk.download('words')
dumpwords(nltk.corpus.words.words(), 'words')
