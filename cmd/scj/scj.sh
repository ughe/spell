#!/usr/bin/env sh

# S. C. Johnson Spell Check implementation by W. Ughetta on 10/2/21
# Case insensitive. Strips non a-z chars. Expects sorted dictionary

## In "Development of a Spelling List" (1982) M. D. McIlroy wrote: "Some
## years ago, S. C. Johnson introduced the UNIX spelling checker spell.
## His idea was simply to look up every word of a document in a standard
## desk dictionary and print a list of the words that were not found.

# Note: command for replacing CRLF with LF in a file:
# awk 'BEGIN {RS="\r\n";ORS="\n"} {print $0}' src > dst

if [ "$#" -ne 2 ]; then echo "usage: $0 dictionary check.txt"; exit 1; fi
if grep -q "$" "$2" || grep -q "$" "$1"; then echo "error: CRLFs found"; exit 1; fi
cat "$2" | sed 's/[^A-Za-z ]//g' | tr ' ' '\n' | sort -fu | comm -23i - "$1"
