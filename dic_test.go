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

func assertGetCharInfo(t *testing.T, cp *charProperty, code_point uint16, default_type uint32, char_type uint32, char_count uint32, group uint32, invoke uint32) {
	v1, v2, v3, v4, v5 := cp.getCharInfo(code_point)
	if v1 != default_type || v2 != char_type || v3 != char_count || v4 != group || v5 != invoke {
		t.Errorf("cp.getCharInfo(%d) is failed", code_point)
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

	assertGetCharInfo(t, cp, 0, 0, 1, 0, 1, 0)        // DEFAULT
	assertGetCharInfo(t, cp, 0x20, 1, 2, 0, 1, 0)     // SPACE
	assertGetCharInfo(t, cp, 0x09, 1, 2, 0, 1, 0)     // SPACE
	assertGetCharInfo(t, cp, 0x6f22, 2, 4, 2, 0, 0)   // KANJI 漢
	assertGetCharInfo(t, cp, 0x3007, 3, 264, 0, 1, 1) // SYMBOL
	assertGetCharInfo(t, cp, 0x31, 4, 16, 0, 1, 1)    // NUMERIC 1
	assertGetCharInfo(t, cp, 0x3042, 6, 64, 2, 1, 0)  // HIRAGANA あ
	assertGetCharInfo(t, cp, 0x4e00, 8, 260, 0, 1, 1) // KANJINUMERIC 一
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
