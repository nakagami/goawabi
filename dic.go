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

func c_str_to_string(data []byte) string {
	i := 0
	for data[i] != 0 {
		i++
	}
	return string(data[:i])
}

// CharProperty

type charProperty struct {
	data           []byte
	category_names []string
	offset         int
}

func newCharProperty(path string) (cp *charProperty, err error) {
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

	cp = new(charProperty)
	cp.data = data
	category_size := int(binary.LittleEndian.Uint32(cp.data))
	cp.category_names = make([]string, category_size)
	for i := 0; i < category_size; i++ {
		cp.category_names[i] = c_str_to_string(cp.data[4+i*32 : 4+(i+1)*32])
	}
	cp.offset = 4 + category_size*32

	return cp, err
}

func (cp *charProperty) getCharInfo(code_point uint16) (uint32, uint32, uint32, uint32, uint32) {
	v := binary.LittleEndian.Uint32(cp.data[cp.offset+int(code_point)*4:])
	default_type := (v >> 18) & 0b11111111
	char_type := v & 0b111111111111111111
	char_count := (v >> 26) & 0b1111
	group := (v >> 30) & 0b1
	invoke := (v >> 31) & 0b1
	return default_type, char_type, char_count, group, invoke
}

// MecabDic

type mecabDic struct {
	data           []byte
	dic_size       int
	lsize          int
	rsize          int
	da_offset      int
	token_offset   int
	feature_offset int
}

func newMecabDic(path string) (m *mecabDic, err error) {
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

	m = new(mecabDic)
	m.data = data
	m.dic_size = int(binary.LittleEndian.Uint32(m.data[0:]) ^ 0xef718f77)
	m.lsize = int(binary.LittleEndian.Uint32(m.data[16:]))
	m.rsize = int(binary.LittleEndian.Uint32(m.data[20:]))
	m.da_offset = 72
	m.token_offset = m.da_offset + int(binary.LittleEndian.Uint32(m.data[24:]))
	m.feature_offset = m.token_offset + int(binary.LittleEndian.Uint32(m.data[28:]))

	return m, err
}

func (m *mecabDic) baseCheck(idx uint32) (int32, uint32) {
	i := m.da_offset + int(idx*8)
	base := int32(binary.LittleEndian.Uint32(m.data[i:]))
	check := binary.LittleEndian.Uint32(m.data[i+4:])

	return base, check
}

func (m *mecabDic) exactMatchSearch(s []byte) int32 {
	var v int32 = -1
	var p uint32
	b, _ := m.baseCheck(0)
	for _, item := range s {
		p = uint32(b+int32(item)) + 1
		base, check := m.baseCheck(p)
		if b == int32(check) {
			b = base
		} else {
			return v
		}
	}

	p = uint32(b)
	n, check := m.baseCheck(p)
	if b == int32(check) && n < 0 {
		v = -n - 1
	}

	return v
}

func (m *mecabDic) commonPrefixSearch(s []byte) [][2]int32 {
	results := make([][2]int32, 0)
	var p uint32
	b, _ := m.baseCheck(0)
	for i, item := range s {
		p = uint32(b)
		n, check := m.baseCheck(p)
		if b == int32(check) && n < 0 {
			results = append(results, [2]int32{-n - 1, int32(i)})
		}
		p = uint32((b + int32(item))) + 1
		base, check := m.baseCheck(p)
		if b == int32(check) {
			b = base
		} else {
			return results
		}
	}
	p = uint32(b)

	n, check := m.baseCheck(p)
	if b == int32(check) && n < 0 {
		results = append(results, [2]int32{-n - 1, int32(len(s))})
	}

	return results
}

// Matrix

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
