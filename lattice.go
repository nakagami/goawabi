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

// Node

type Node struct {
	entry      *DicEntry
	pos        int32
	epos       int32
	index      int32
	left_id    int32
	right_id   int32
	cost       int32
	min_cost   int32
	back_pos   int32
	back_index int32
}

func newBos() *Node {
	node := new(Node)
	node.entry = nil
	node.pos = 0
	node.epos = 1
	node.index = 0
	node.left_id = -1
	node.right_id = 0
	node.cost = 0
	node.min_cost = 0
	node.back_pos = -1
	node.back_index = -1

	return node
}

func newEos(pos int32) *Node {
	node := new(Node)
	node.entry = nil
	node.pos = pos
	node.epos = pos + 1
	node.index = 0
	node.left_id = 0
	node.right_id = -1
	node.cost = 0
	node.min_cost = 0x7FFFFFFF
	node.back_pos = -1
	node.back_index = -1

	return node
}

func newNode(e *DicEntry) *Node {
	node := new(Node)
	node.entry = e
	node.pos = 0
	node.epos = 0
	node.index = int32(e.posid)
	node.left_id = int32(e.lc_attr)
	node.right_id = int32(e.rc_attr)
	node.cost = int32(e.wcost)
	node.min_cost = 0x7FFFFFFF
	node.back_pos = -1
	node.back_index = -1

	return node
}

func (node *Node) isBos() bool {
	return node.entry == nil && node.pos == 0
}

func (node *Node) isEos() bool {
	return node.entry == nil && node.pos != 0
}

func (node *Node) nodeLen() int {
	if node.entry == nil {
		return len(node.entry.original)
	}

	return 1 // BOS or EOS
}

// Lattice

type Lattice struct {
	snodes [][]Node
	enodes [][]Node
	p      int32
}
