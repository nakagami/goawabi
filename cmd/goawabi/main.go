package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"github.com/nakagami/goawabi"
)

func main() {
	tokenizer, err := goawabi.NewTokenizer("")
	if err != nil {
		panic(err)
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	morphemes, err := tokenizer.Tokenize(strings.TrimSpace(string(input)))
	if err != nil {
		panic(err)
	}

	for _, m := range morphemes {
		fmt.Printf("%s\t%s\n", m[0], m[1])
	}
	fmt.Printf("EOS\n")
}
