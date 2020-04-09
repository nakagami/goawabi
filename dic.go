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

func utf8ToUcs2(s []byte, index int) (uint16, int) {
	// utf8 to ucs2(16bit) code and it's array size
	ln := 0

	if (s[index] & 0b10000000) == 0b00000000 {
		ln = 1
	} else if (s[index] & 0b11100000) == 0b11000000 {
		ln = 2
	} else if (s[index] & 0b11110000) == 0b11100000 {
		ln = 3
	} else if (s[index] & 0b11111000) == 0b11110000 {
		ln = 4
	}

	var ch32 uint32
	switch ln {
	case 1:
		ch32 = uint32(s[index+0])
	case 2:
		ch32 = uint32(s[index+0]&0x1F) << 6
		ch32 |= uint32(s[index+1] & 0x3F)
	case 3:
		ch32 = uint32(s[index+0]&0x0F) << 12
		ch32 |= uint32(s[index+1]&0x3F) << 6
		ch32 |= uint32(s[index+2] & 0x3F)
	case 4:
		ch32 = uint32(s[index+0]&0x07) << 18
		ch32 |= uint32(s[index+1]&0x3F) << 12
		ch32 |= uint32(s[index+2]&0x3F) << 6
		ch32 |= uint32(s[index+3] & 0x03F)
	}

	// ucs4 to ucs2
	var ch16 uint16
	if ch32 < 0x10000 {
		ch16 = uint16(ch32)
	} else {
		ch16 = uint16((((ch32-0x10000)/0x400 + 0xD800) << 8) + ((ch32-0x10000)%0x400 + 0xDC00))
	}
	return ch16, ln
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

func (cp *charProperty) getGroupLength(s []byte, default_type uint32, max_count int) int {
	var i, char_count int

	for i < len(s) {
		ch16, ln := utf8ToUcs2(s, i)
		_, t, _, _, _ := cp.getCharInfo(ch16)

		if ((1 << default_type) & t) != 0 {
			i += ln
			char_count += 1
			if max_count != 0 && max_count == char_count {
				break
			}
		} else {
			break
		}
	}
	return i
}

func (cp *charProperty) getCountLength(s []byte, count int) int {
	var i int

	for j := 0; j < count; j++ {
		_, ln := utf8ToUcs2(s, i)
		i += ln
	}
	return i
}

func (cp *charProperty) getUnknownLengths(s []byte) (uint32, []int, bool) {
	// get unknown word bytes length vector
	ln_list := make([]int, 0)
	ch16, _ := utf8ToUcs2(s, 0)
	default_type, _, count, group, invoke := cp.getCharInfo(ch16)
	if group != 0 {
		ln_list = append(ln_list, cp.getGroupLength(s, default_type, int(count)))
	} else {
		ln_list = append(ln_list, cp.getCountLength(s, int(count)))
	}

	return default_type, ln_list, invoke == 1
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

func (m *mecabDic) getEntriesByIndex(idx int, count int, s string) []*DicEntry {
	results := make([]*DicEntry, 0)
	for i := 0; i < count; i++ {
		d := new(DicEntry)
		offset := m.token_offset + (idx+i)*16
		d.original = s
		d.lc_attr = binary.LittleEndian.Uint16(m.data[offset:])
		d.rc_attr = binary.LittleEndian.Uint16(m.data[offset+2:])
		d.posid = binary.LittleEndian.Uint16(m.data[offset+4:])
		d.wcost = int16(binary.LittleEndian.Uint16(m.data[offset+6:]))
		feature := int(binary.LittleEndian.Uint32(m.data[offset+8:]))
		d.feature = c_str_to_string(m.data[m.feature_offset+feature:])
		results = append(results, d)
	}

	return results
}

func (m *mecabDic) getEntries(result int, s string) []*DicEntry {
	return m.getEntriesByIndex(result>>8, result&0xff, s)
}

func (m *mecabDic) lookup(s []byte) []*DicEntry {
	results := make([]*DicEntry, 0)

	for _, v := range m.commonPrefixSearch(s) {
		result, ln := v[0], v[1]
		index := int(result >> 8)
		count := int(result & 0xff)
		newResults := m.getEntriesByIndex(index, count, string(s[:ln]))
		results = append(results, newResults...)
	}

	return results
}

func (m *mecabDic) lookupUnknowns(s []byte, cp *charProperty) ([]*DicEntry, bool) {
	default_type, ln_list, invoke := cp.getUnknownLengths(s)
	index := m.exactMatchSearch([]byte(cp.category_names[int(default_type)]))
	results := make([]*DicEntry, 0)
	for _, ln := range ln_list {
		newResults := m.getEntries(int(index), string(s[:ln]))
		results = append(results, newResults...)
	}
	return results, invoke
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

func (m *matrix) getTransCost(id1 int, id2 int) int32 {
	i := (id2*m.lsize+id1)*2 + 4
	return int32(int16(binary.LittleEndian.Uint16(m.data[i:])))
}
