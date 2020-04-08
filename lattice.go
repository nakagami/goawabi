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

func (node *Node) nodeLen() int32 {
	if node.entry == nil {
		return int32(len(node.entry.original))
	}

	return 1 // BOS or EOS
}

// Lattice

type Lattice struct {
	snodes [][]*Node
	enodes [][]*Node
	p      int32
}

func newLattice(size int) (lat *Lattice, err error) {
	lat = new(Lattice)
	lat.snodes = make([][]*Node, size+2)
	lat.enodes = make([][]*Node, size+2)
	for i := 0; i < size+2; i++ {
		lat.snodes[i] = make([]*Node, 0)
		lat.enodes[i] = make([]*Node, 0)
	}
	bos := newBos()
	lat.snodes[0] = append(lat.snodes[0], bos)
	lat.enodes[1] = append(lat.snodes[1], bos)
	lat.p = 1

	return lat, err
}

func (lat *Lattice) add(node *Node, m *matrix) {
	min_cost := node.min_cost
	best_node := lat.enodes[lat.p][0]

	for _, enode := range lat.enodes[lat.p] {
		cost := enode.min_cost + m.getTransCost(int(enode.right_id), int(enode.left_id))
		if cost < min_cost {
			min_cost = cost
			best_node = enode
		}
	}

	node.min_cost = min_cost + node.cost
	node.back_index = best_node.index
	node.back_pos = best_node.pos
	node.pos = lat.p
	node.epos = lat.p + node.nodeLen()

	node.index = int32(len(lat.snodes[lat.p]))

	node_pos := node.pos
	node_epos := node.epos
	lat.snodes[node_pos] = append(lat.snodes[node_pos], node)
	lat.enodes[node_epos] = append(lat.enodes[node_epos], node)
}

func (lat *Lattice) forward() int {
	old_p := lat.p
	lat.p += 1
	for len(lat.enodes[lat.p]) == 0 {
		lat.p += 1
	}
	return int(lat.p - old_p)
}

/*
func (lat *Lattice) end(m *matrix) int {
}

func (lat *Lattice) backward() []*Node {
}
*/
