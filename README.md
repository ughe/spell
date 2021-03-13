# spell

Transliteration of [v10spell](https://github.com/arnoldrobbins/v10spell) from C to Go.

> WIP. Not functioning yet.

## Example

```
go install ./...
cd dictionaries
pcode list british local stop  > brspell
pcode list american local stop > amspell
```

Outputs `brspell` and `amspell` along with `stderr`:

```
words = 31287; codes = 285
output bytes = 164027
words = 31292; codes = 284
output bytes = 163957
```
