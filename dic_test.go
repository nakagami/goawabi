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

func TestMatrix(t *testing.T) {
	mecabrc_map, _ := get_mecabrc_map("")
	path := get_dic_path(mecabrc_map, "matrix.bin")
	m, err := newMatrix(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.getTransCost(555, 1283) != 340 {
		t.Errorf("getTransCost(555, 1283)")
	}
	if m.getTransCost(10, 1293) != -1376 {
		t.Errorf("getTransCost(10, 1293)")
	}
}

func TestCharPropery(t *testing.T) {
	mecabrc_map, _ := get_mecabrc_map("")
	path := get_dic_path(mecabrc_map, "char.bin")
	cp, err := newCharProperty(path)
	if err != nil {
		t.Fatal(err)
	}

	for i, s := range []string{"DEFAULT", "SPACE", "KANJI", "SYMBOL", "NUMERIC", "ALPHA", "HIRAGANA", "KATAKANA", "KANJINUMERIC", "GREEK", "CYRILLIC"} {
		if s != cp.category_names[i] {
			t.Errorf("category_name %s is invalid", cp.category_names[i])
		}
	}

}

func TestMecabDic(t *testing.T) {
	mecabrc_map, _ := get_mecabrc_map("")
	path := get_dic_path(mecabrc_map, "sys.dic")
	m, err := newMecabDic(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.dic_size != 49199027 {
		t.Errorf("sys.dic is incollect size %d", m.dic_size)
	}
}
