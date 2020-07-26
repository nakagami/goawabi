package main

import (
	"flag"
	"fmt"
	"github.com/nakagami/goawabi"
	"io/ioutil"
	"os"
	"regexp"
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
	flag.Parse()

	tokenizer, err := goawabi.NewTokenizer("")
	if err != nil {
		panic(err)
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	for _, s := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(input), -1) {

		if *n > 1 {
			morphemes_list, err := tokenizer.TokenizeNBest(s, *n)
			if err != nil {
				panic(err)
			}
			for _, m := range morphemes_list {
				print(m)
			}

		} else {
			morphemes, err := tokenizer.Tokenize(s)
			if err != nil {
				panic(err)
			}
			print(morphemes)
		}

	}

}
