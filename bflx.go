// Copyright 2017 Joubin Houshyar
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package bflx

import (
	"fmt"
)

// ----------------------------------------------------------------------
// interpreter memory object (level)
// ----------------------------------------------------------------------

// level data cell
type memobj struct {
	data []byte
	dx   int
}

// returns pointer to new instance of memobj
func newMemobj() *memobj {
	return &memobj{
		data: make([]byte, 1),
	}
}

// moves data cursor forward by 1.
// if index exceeds capacity, capacity is increased.
func (p *memobj) forward() {
	if p.dx == len(p.data)-1 {
		var b byte
		p.data = append(p.data, b)
	}
	p.dx++
	fmt.Printf("debug - > - %d len:%d\n", p.dx, len(p.data))
}

// move data cursor back by 1.
// if index underflows (>0) move to end per circular buffer semantics.
func (p *memobj) back() {
	if p.dx == 0 {
		p.dx = len(p.data)
	}
	p.dx--
	fmt.Printf("debug - < - %d\n", p.dx)
}

// decrement current cell value
func (p *memobj) decrement() {
	p.data[p.dx]--
}

// increment current cell value
func (p *memobj) increment() {
	p.data[p.dx]++
}

// invert current cell bits
func (p *memobj) invert() {
	p.data[p.dx] ^= 0xff
}

// returns value of current cell
func (p *memobj) Get() byte {
	return p.data[p.dx]
}

// sets value of current cell
func (p *memobj) Set(b byte) {
	p.data[p.dx] = b
}

// ----------------------------------------------------------------------
// interpreter
// ----------------------------------------------------------------------

// type wrapper for interpreter state
type interpreter struct {
	register [16]byte  // indexed & special registers
	rx       int       // register index
	level    []*memobj // level data
	lx       int       // level index
}

// returns pointer to new instance of a BFLX interpreter
func NewInterpreter() *interpreter {
	p := &interpreter{}
	p.level = append(p.level, newMemobj())
	return p
}

// increment level counter
// if overflow, allocate new data level
func (p *interpreter) levelUp() {
	if p.lx == len(p.level)-1 {
		p.level = append(p.level, newMemobj())
	}
	p.lx++
}

// decrement level counter
// if underflow, go to top.
func (p *interpreter) levelDown() {
	if p.lx == 0 {
		p.lx = len(p.level)
	}
	p.lx--
}

// go to top level
func (p *interpreter) levelTop() {
	p.lx = len(p.level) - 1
}

// go to bottom level
func (p *interpreter) levelFloor() {
	p.lx = 0
}

// interpreter run loop.
func (p *interpreter) Run(program string) string {
	var out []byte
	var inst = []byte(program)
	for ix := 0; ix < len(inst); ix++ {
		d := 1
		fmt.Printf("debug - token:%c - rx:%d\n", inst[ix], p.rx)
		switch {
		case inst[ix] == '[' && p.level[p.lx].Get() == 0:
			for d > 0 {
				ix++
				switch inst[ix] {
				case '[':
					d++
				case ']':
					d--
				}
			}
		case inst[ix] == ']' && p.level[p.lx].Get() != 0:
			for d > 0 {
				ix--
				switch inst[ix] {
				case '[':
					d--
				case ']':
					d++
				}
			}
		case inst[ix] >= '0' && inst[ix] <= '9':
			p.rx = int(inst[ix] - 48)
			fmt.Printf("debug - register[%d]=%d\n", p.rx, p.register[p.rx])
		case inst[ix] == '#':
			p.register[p.rx] = p.level[p.lx].Get()
			fmt.Printf("debug - register[%d]=%d\n", p.rx, p.register[p.rx])
		case inst[ix] == '%':
			fmt.Printf("debug - register[%d]=%d level:%d\n", p.rx, p.register[p.rx], p.lx)
			p.level[p.lx].Set(p.register[p.rx])
		case inst[ix] == '+':
			p.level[p.lx].increment()
		case inst[ix] == '-':
			p.level[p.lx].decrement()
		case inst[ix] == '>':
			p.level[p.lx].forward()
		case inst[ix] == '<':
			p.level[p.lx].back()
		case inst[ix] == '(':
			p.level[p.lx].dx = 0
		case inst[ix] == ')':
			p.level[p.lx].dx = len(p.level[p.lx].data) - 1
		case inst[ix] == '^':
			p.levelUp()
		case inst[ix] == 'v':
			p.levelDown()
		case inst[ix] == 'T':
			p.levelTop()
		case inst[ix] == '_':
			p.levelFloor()
		case inst[ix] == 'w':
			out = append(out, p.level[p.lx].Get())
			p.level[p.lx].forward()
		case inst[ix] == 'n':
			numrep := fmt.Sprintf("%d", p.level[p.lx].Get())
			out = append(out, []byte(numrep)...)
			p.level[p.lx].forward()
		case inst[ix] == 'N':
			numrep := fmt.Sprintf("%03d", p.level[p.lx].Get())
			out = append(out, []byte(numrep)...)
			p.level[p.lx].forward()
		case inst[ix] == '?':
			var b byte
			fmt.Scanf("%c\n", &b)
			fmt.Printf("debug input:%d\n", b)
			p.level[p.lx].Set(b)
			p.level[p.lx].forward()
		default:
		}
	}

	return string(out)
}
