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
	"os"
	"testing"
)

func TestFindMecabRc(t *testing.T) {
	_, err := find_mecabrc()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMecabRcMap(t *testing.T) {
	mecabrc_map, err := get_mecabrc_map("")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := mecabrc_map["dicdir"]; !ok {
		t.Fatalf("Can't find dicdir")
	}
}

func TestGetDicPath(t *testing.T) {
	mecabrc_map, _ := get_mecabrc_map("")
	for _, s := range []string{"sys.dic", "unk.dic", "matrix.bin", "char.bin"} {
		path := get_dic_path(mecabrc_map, s)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			t.Fatalf("Can't find %s", path)
		}
	}
}
