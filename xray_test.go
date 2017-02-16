package main

import (
	"bytes"
	"testing"
	"bufio"
	"io"
)

func b(s string) []byte{
	return []byte(s)
}

func TestParseRun(t *testing.T) {
	type Test struct {
		Name  string
		Input []byte
		Want  []Item
	}
	mk := func(name string, input string, it ...Item) Test{
		return Test{
			Name: name, 
			Input: []byte(input),
			Want: it,
		}
	}
	
	var Tests = []Test{
		mk("ParseRunABC", "ABC",  &Run{v: b("ABC"), items: []Item{}}, EOF{io.EOF}),
	}

	for _, v := range Tests {
		p := &parser{br: bufio.NewReader(bytes.NewReader([]byte(v.Input)))}
		items := p.parse()
		for i := range items{
			if !bytes.Equal(v.Want[i].Bytes(), items[i].Bytes()) {
				t.Fail()
			}		
		}

	}

}
