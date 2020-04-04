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
	"encoding/binary"
	"os"
	"syscall"
)

type DicEntry struct {
	original string
	lc_attr  uint16
	rc_attr  uint16
	posid    uint16
	wcost    int16
	feature  string
}

type charProperty struct {
	data           []byte
	category_names []string
	offset         int
}

type mecabDic struct {
	data           []byte
	dic_size       int
	lsize          int
	rsize          int
	da_offset      int
	token_offset   int
	feature_offset int
}

type matrix struct {
	data  []byte
	lsize int
	rsize int
}

func newMatrix(path string) (m *matrix, err error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()
	finfo, _ := f.Stat()
	data, err := syscall.Mmap(int(f.Fd()), 0, int(finfo.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	m = new(matrix)
	m.data = data
	m.lsize = int(binary.LittleEndian.Uint16(m.data))
	m.rsize = int(binary.LittleEndian.Uint16(m.data[2:]))

	return m, err
}

func (m *matrix) getTransCost(id1 int, id2 int) int {
	i := (id2*m.lsize+id1)*2 + 4
	return int(int16(binary.LittleEndian.Uint16(m.data[i:])))
}
