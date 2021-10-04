# Generate dictionary (word list) using WordNet
import nltk
from nltk.corpus import wordnet as wn

def isword(word):
    return all([x >= 'a' and x <= 'z' for x in word.lower()])

nltk.download('wordnet')

names = [str(x) for x in wn.all_lemma_names()]
words = sorted(filter(lambda x: x, [x if isword(x) else None for x in names]), key=str.casefold)
other = sorted(filter(lambda x: x, [x if not isword(x) else None for x in names]), key=str.casefold)

with open('dictionary.txt', 'w') as f:
    f.write('\n'.join(words))

with open('wordnet.txt', 'w') as f:
    f.write('\n'.join(words))

with open('othernet.txt', 'w') as f:
    f.write('\n'.join(other))

print('Success. Created wordnet.txt othernet.txt.\nTotal: %d. Words: %d. Other: %d.' % (len(names), len(words), len(other)))
