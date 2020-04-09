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

type Tokenizer struct {
	sys_dic  *mecabDic
	user_dic *mecabDic
	cp       *charProperty
	unk_dic  *mecabDic
	m        *matrix
}

func NewTokenizer(path string) (*Tokenizer, error) {
	tok := new(Tokenizer)
	mecabrc_map, _ := get_mecabrc_map(path)
	sys_dic, err := newMecabDic(get_dic_path(mecabrc_map, "sys.dic"))
	if err != nil {
		return tok, err
	}
	tok.sys_dic = sys_dic

	if val, ok := mecabrc_map["userdic"]; ok {
		user_dic, err := newMecabDic(val)
		if err != nil {
			return tok, err
		}
		tok.user_dic = user_dic
	}
	cp, err := newCharProperty(get_dic_path(mecabrc_map, "char.bin"))
	if err != nil {
		return tok, err
	}

	tok.cp = cp

	unk_dic, err := newMecabDic(get_dic_path(mecabrc_map, "unk.dic"))
	if err != nil {
		return tok, err
	}
	tok.unk_dic = unk_dic
	m, err := newMatrix(get_dic_path(mecabrc_map, "matrix.bin"))
	if err != nil {
		return tok, err
	}
	tok.m = m

	return tok, err
}

func (tok *Tokenizer) buildLattice(str string) (*Lattice, error) {
	s := []byte(str)
	lat, err := newLattice(len(s))
	pos := 0
	for pos < len(s) {
		matched := false

		// user_dic
		if tok.user_dic != nil {
			user_entries := tok.user_dic.lookup(s[pos:])
			if len(user_entries) > 0 {
				for _, entry := range user_entries {
					lat.add(newNode(entry), tok.m)
				}
				matched = true
			}
		}

		// sys_dic
		sys_entries := tok.sys_dic.lookup(s[pos:])
		if len(sys_entries) > 0 {
			for _, entry := range sys_entries {
				lat.add(newNode(entry), tok.m)
			}
			matched = true
		}

		// unknown
		unk_entries, invoke := tok.unk_dic.lookupUnknowns(s[pos:], tok.cp)
		if invoke || !matched {
			for _, entry := range unk_entries {
				lat.add(newNode(entry), tok.m)
			}

		}

		pos += lat.forward()
	}

	lat.end(tok.m)
	return lat, err
}

func (tok *Tokenizer) Tokenize(str string) ([][2]string, error) {
	lat, err := tok.buildLattice(str)
	if err != nil {
		return nil, err
	}
	nodes := lat.backward()

	morphemes := make([][2]string, 0)

	for i := 1; i < len(nodes)-1; i++ {
		morphemes = append(morphemes, [2]string{nodes[i].entry.original, nodes[i].entry.feature})
	}

	return morphemes, nil
}
