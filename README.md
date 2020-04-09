# goawabi

`goawabi` is a morphological analyzer using mecab dictionary, written in Go.

See also python wrapper `awabi` https://github.com/nakagami/awabi .

## Requirements and how to install

MeCab https://taku910.github.io/mecab/ and related dictionary is required.

### Debian/Ubuntu
```
$ sudo apt install mecab
$ go get github.com/nakagami/goawabi/cmd/goawabi
```

### Mac OS X (homebrew)
```
$ brew install mecab
$ brew install mecab-ipadic
$ go get github.com/nakagami/goawabi/cmd/goawabi
```

