## Dictionary

Word lists contain one word per line with only ascii letters `A` to `Z`, upper or lower case. The word list is sometimes found in `/usr/share/dict/words`. Here are a couple of word lists:

| Filename | Line Count | Description |
| --- | --- | --- |
| `linuxwords` | 45,402 | Abridged list of linux words. From [duke.edu](https://www.cs.duke.edu/~ola/ap/linuxwords). Cited by [wikipedia](https://en.wikipedia.org/wiki/Words_(Unix)) |
| `web2` | 234,936 | Original linux words. From [netbsd.org](http://cvsweb.netbsd.org/bsdweb.cgi/~checkout~/src/share/dict/web2?rev=1.1&content-type=text/plain&only_with_tag=MAIN). Cited by [stackexchange](https://unix.stackexchange.com/a/286790) |
| `words` | 236,734 | Linux words from Python NLTK. See `./nltk-gen.py` |
| `wordnet.txt` | 77,503 | WordNet words from Python NLTK. See `./nltk-gen.py` |

## web2

The following quote is from [netbsd.org](http://cvsweb.netbsd.org/bsdweb.cgi/src/share/dict/README?rev=1.1&content-type=text/x-cvsweb-markup&only_with_tag=MAIN) and referenced from [stackexchange](https://unix.stackexchange.com/a/286790):

> Welcome to web2 (Webster's Second International) all 234,936 words worth.
> The 1934 copyright has elapsed, according to the supplier.  The
> supplemental 'web2a' list contains hyphenated terms as well as assorted
> noun and adverbial phrases.  The wordlist makes a dandy 'grep' victim.
>
>      -- James A. Woods    {ihnp4,hplabs}!ames!jaw    (or jaw@riacs)
