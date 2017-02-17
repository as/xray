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

func TestParseRunData(t *testing.T){
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
		mk("ParseASCIIABC", "ABC", &Run{ v: b("ABC"), 
			items: []Item{&Run{v: b("ABC"), items: []Item{&ASCII{v: b("ABC")}}}}}, End),
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
		p.combineRuns(&items)
		p.typeCheckAll(&items)
		check(t, v.Name, v.Want, items)

	}
}

func check(t *testing.T, name string, X, Y []Item){
		if len(Y) != len(X){
			t.Logf("%s\nitem count differs\nhave: %d\nwant: %d\n", name, len(Y),len(X))
			t.Fail()
			return		
		}
		for i := range X{
			want, have := X[i], Y[i]
			switch yt := have.(type){
			case *Run:
				if xt, ok := want.(*Run); ok{
					check(t, name, xt.items, yt.items)
				} else {
					t.Logf("%s\nhave: %T\nwant: %T\n", name, have, want)
					t.Fail()				
				}
			case Item:
				if !bytes.Equal(want.Bytes(), have.Bytes()) {
					t.Logf("%s\nhave: %#v\nwant: %#v\n", name, have, want)
					t.Fail()
				}
			}		
		}
}