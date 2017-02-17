package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type parser struct {
	br           *bufio.Reader
	tok, last    byte
	n            int
	err, lasterr error
	indent       int
}

func (p *parser) next() byte {
	p.last = p.tok
	p.lasterr = p.err
	p.tok, p.err = p.br.ReadByte()
	p.printTrace(fmt.Sprintf("next: %x", p.last))
	return p.last
}

func (p *parser) parse() (items []Item) {
	p.next()
	for {
		it := p.parseAny()
		items = append(items, it)
		if _, ok := it.(EOF); ok {
			break
		}
	}
	return
}
func (p *parser) parseAny() (item Item) {
	defer un(trace(p, "parseAny"))
	p.next()
	if p.lasterr != nil {
		return EOF{p.lasterr}
	}
	if p.tok == p.last {
		rep := p.parseRepeat()
		if rep.n < 8{
			item = &Run{
				v: bytes.Repeat([]byte{rep.b}, rep.n),
			}
		} else {
			item = rep
		}
	} else {
		item = p.parseRun()
	}
	return
}

func (p *parser) parseRepeat() (rep *Repeat) {
	defer un(trace(p, "parseRepeat"))
	rep = &Repeat{p.last, 1}
	for p.last == p.tok && p.lasterr == nil {
		p.next()
		rep.n++
	}
	return
}

func (p *parser) parseRun() (run *Run) {
	defer un(trace(p, "parseRun"))
	defer func() { fmt.Printf("parseRun: %#v\n", run) }()
	run = &Run{
		v: []byte{},
		items: []Item{},
	}
	for p.last != p.tok && p.lasterr == nil {
		run.v = append(run.v, p.last)
		p.next()
	}
	if p.err == nil {
		run.v = append(run.v, p.last)
	}


	return
}

func (p *parser) typeCheck(it Item){
	switch t := it.(type){
	case *Run:
		p.parseRunData(t)
		p.parseNUL(t)
		p.parseConformData(t)
	}
}

// parseFunc parses the input by calling fn for every consecutive byte
// sequence until fn returns false
func (p *parser) parseFunc(b []byte, fn func(byte) bool) []byte {
	buf := new(bytes.Buffer)
	for _, x := range b {
		if !fn(x) {
			break
		}
		buf.WriteByte(x)
	}
	return buf.Bytes()
}

func (p *parser) parseNUL(run *Run) {
	var (
		i  = 0
		it Item
	)
	peek := func() Item {
		if i+1 < len(run.items) {
			return run.items[i+1]
		}
		return EOF{io.EOF}
	}
	for i := 0; i < len(run.items); i++{
		it = run.items[i]
		next := peek()
		switch t := it.(type) {
		case *ASCII:
			if !p.acceptNUL(next) {
				break
			}
			t.v = append(t.v, 0)
			t.null = true
			run.Delete(i+1)
		case *UTF16:
			if !p.acceptNUL(t) {
				break
			}
			i++
			if !p.acceptNUL(t) {
				i--
				break
			}
			t.v = append(t.v, 0, 0)
			run.Delete(i+1)
		}
	}
}

func (p *parser) parseConformData(run *Run) {
	lastnumber := false
	var prev Item
	for i, v := range run.items {	
		switch t := v.(type) {
		case *ASCII:
			if lastnumber {
				n := int(prev.(*Number).v)
				if lenlen, len := conformsAny(t, n); lenlen != 0 {
					run.items[i-1] = &Conform{lenlen, len, v}
					run.Delete(i)					
				}
			}
			lastnumber = false
		case *UTF16:
			if lastnumber {
				n := int(prev.(*Number).v) * 2
				if lenlen, len := conformsAny(t, n); lenlen != 0 {
					run.items[i-1] = &Conform{lenlen, len, v}
					run.Delete(i)
				}
			}
			lastnumber = false
		case *Number:
			lastnumber = true
		}
		prev = v
	}
}

func (p *parser) parseRunData(run *Run) {
	defer un(trace(p, "parseRunData"))
	b := run.v

	add := func(it Item) {
		run.items = append(run.items, it)
	}
	debugadd := func(s string, it Item){
		fmt.Printf("add item (%s): %#v\n", s, it)
		add(it)
	}
	
	for len(b) > 0 {
		numeric := p.parseFunc(b, func(b byte) bool { return !common(b) && b != 0 })
		advance := len(numeric)
		if advance > len(b) {
			b = b[len(b)-1:]
		} else {
			b = b[len(numeric):]
		}
		for _, x := range numeric {
			// to make things work, convert every number to be 1 byte wide
			// future: size appropriately with multiple passes
			debugadd("numeric",&Number{int64(x), 1})
		}
		utfs := p.parseUTF16(b)
		ascii := p.parseASCII(b)
		if n := len(utfs.v) * 2; n > 0 {
			debugadd("utf16",utfs)
			b = b[n:]
		} else if n := len(ascii.v); n > 0 {
			debugadd("ascii",ascii)
			b = b[n:]
		} else {
			// numeric function didn't parse this so just make progress for now
			// possibly a stray null terminator
			if len(b) > 0 {
				debugadd("stray", &Number{int64(b[0]), 1})
				b = b[1:]
			}
		}
	}
	return
}

func (p *parser) parseUTF16(b []byte) *UTF16 {
	defer un(trace(p, "parseUTF16"))
	br := new(bytes.Buffer)
	i := 0
	null := false
	for ; i+1 < len(b); i += 2 {
		if common(b[i]) && b[i+1] == 0 {
			br.WriteByte(b[i])
		} else {
			break
		}
	}
	if i+2 < len(b) && b[i+2] == 0 && b[i+1] == 0 {
		null = true
		br.Write([]byte{0, 0})
	}
	if i == 0{
		return &UTF16{}
	}
	fmt.Printf("i is %d and len(br) is %d\n", i, len(br.Bytes()))
	return &UTF16{br.Bytes(), null}
}

func (p *parser) parseASCII(b []byte) *ASCII {
	defer un(trace(p, "parseASCII"))
	br := new(bytes.Buffer)
	i := 0
	null := false
	for ; i < len(b); i++ {
		if common(b[i]) {
			br.WriteByte(b[i])
		} else {
			break
		}
	}
	if i+1 < len(b) && b[i+1] == 0 {
		null = true
		br.WriteByte(0)
	}
	if i == 0{
		return &ASCII{}
	}
	return &ASCII{v: br.Bytes(), null: null}
}

func (p *parser) acceptNUL(it Item) bool {
	if num, ok := it.(*Number); !ok {
		return false
	} else if num.v != 0 {
		return false
	}
	return true
}
