# goawabi

`goawabi` is a morphological analyzer using mecab dictionary, written in Go.

See also an original Rust implementation `awabi` https://github.com/nakagami/awabi .

## Requirements and how to install

### Install Mecab dictionary data

MeCab https://taku910.github.io/mecab/ and related dictionary is required.

It is recommended to install mecab to update the dictionary and check its operation.

#### Debian/Ubuntu
```
$ sudo apt install mecab
$ sudo apt install mecab-ipadic-utf8
```

#### Mac OS X (homebrew)
```
$ brew install mecab
$ brew install mecab-ipadic
```

## How to use

Goawabi can execute as a command or called from a library

### goawabi command

#### Install goawabi command

```
$ go install github.com/nakagami/goawabi/cmd/goawabi@latest
```

#### Command Excample

```
$ echo 'すもももももももものうち' |goawabi
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
の	助詞,連体化,*,*,*,*,の,ノ,ノ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
EOS
```

```
$ echo 'すもももももももものうち' |goawabi  -N 2
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
の	助詞,連体化,*,*,*,*,の,ノ,ノ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
EOS
すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
もも	名詞,一般,*,*,*,*,もも,モモ,モモ
も	助詞,係助詞,*,*,*,*,も,モ,モ
の	助詞,連体化,*,*,*,*,の,ノ,ノ
うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
EOS
```

### use as library

See main as sample code.

- tokensize https://github.com/nakagami/goawabi/blob/master/cmd/goawabi/main.go#L48
- N best match https://github.com/nakagami/goawabi/blob/master/cmd/goawabi/main.go#L39

## See also

- awabi https://github.com/nakagami/awabi Rust implementation
- pyawabi https://github.com/nakagami/pyawabi Python wrapper
