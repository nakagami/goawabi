package main

import (
	"flag"
	"fmt"
	"github.com/nakagami/goawabi"
	"io/ioutil"
	"os"
	"strings"
)

func print(morphemes [][2]string) {
	for _, m := range morphemes {
		fmt.Printf("%s\t%s\n", m[0], m[1])
	}
	fmt.Printf("EOS\n")
}

func main() {
	var (
		n = flag.Int("N", 1, "N best")
	)

	tokenizer, err := goawabi.NewTokenizer("")
	if err != nil {
		panic(err)
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if *n > 1 {
		morphemes_list, err := tokenizer.TokenizeNBest(strings.TrimSpace(string(input)), *n)
		if err != nil {
			panic(err)
		}
		for _, m := range morphemes_list {
			print(m)
		}

	} else {
		morphemes, err := tokenizer.Tokenize(strings.TrimSpace(string(input)))
		if err != nil {
			panic(err)
		}
		print(morphemes)
	}
}
