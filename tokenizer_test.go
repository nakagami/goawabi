/*****************************************************************************
MIT License

Copyright (c) 2020 Hajime Nakagami

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*****************************************************************************/

package goawabi

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	results1 := [][2]string{
		{"すもも", "名詞,一般,*,*,*,*,すもも,スモモ,スモモ"},
		{"も", "助詞,係助詞,*,*,*,*,も,モ,モ"},
		{"もも", "名詞,一般,*,*,*,*,もも,モモ,モモ"},
		{"も", "助詞,係助詞,*,*,*,*,も,モ,モ"},
		{"もも", "名詞,一般,*,*,*,*,もも,モモ,モモ"},
		{"の", "助詞,連体化,*,*,*,*,の,ノ,ノ"},
		{"うち", "名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ"},
	}

	tokenizer, err := NewTokenizer("")
	if err != nil {
		t.Fatal(err)
	}

	morphemes, err := tokenizer.Tokenize("すもももももももものうち")
	if err != nil {
		t.Fatal(err)
	}
	for i, m := range morphemes {
		if results1[i][0] != m[0] || results1[i][1] != m[1] {
			t.Errorf("Tokenize() failed:%s,%s", m[0], m[1])
		}
	}
}
