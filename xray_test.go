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
	End := EOF{io.EOF}
	Empty := []Item{}
	var Tests = []Test{
		mk("ParseRunABC", "ABC",
			&Run{v: b("ABC"), items: Empty}, End),
		mk("ParseRunAAA", "AAA",
			&Run{v: b("AAA"), items: Empty}, End),
		mk("ParseRepeatA", "AAAAAAAA",
			&Repeat{b: 'A', n: 8},  End),
		mk("ParseNonRepeatA", "AAAAAAA",
			&Run{v: b("AAAAAAA"), items: Empty}, End),
		mk("ParseRunRepeat", "AAAABBBBBBBBBB", 
			&Run{v: b("AAAA"), items: Empty},
			&Repeat{b: 'B', n: 10},
			End),
	}

	for _, v := range Tests {
		p := &parser{br: bufio.NewReader(bytes.NewReader([]byte(v.Input)))}
		items := p.parse()
		if len(v.Want) != len(items){
			t.Logf("%s\nitem count differs\nhave: %d\nwant: %d\n", v.Name, len(v.Want),len(items))
			t.Fail()		
		}
		for i := range items{
			want, have := v.Want[i].Bytes(), items[i].Bytes()
			if !bytes.Equal(want, have) {
				t.Logf("%s\nhave: %q\nwant: %q\n", v.Name, have, want)
				t.Fail()
			}		
		}

	}

}
