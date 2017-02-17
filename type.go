package main
import (
	"fmt"
	"bytes"
	"github.com/as/bo"
)

type Item interface{
	Bytes() []byte
}

type EOF struct {
	err error
}
type Repeat struct {
	b byte
	n int
}

type Run struct {
	v     []byte
	items []Item
	len int
}

func (r *Run) Delete(i int) Item{
	l := len(r.items)
	if i >= l || i < 0 || l == 0{
		return nil
	}
	it := r.items[i]
	if i+1 < l{
		copy(r.items[i:], r.items[i+1:])
	}
	r.items[l-1] = EOF{nil}
	return it
}

type Number struct {
	v     int64
	width int
}

// Conform is a box for length-prefixed items
// such as //wire9 n[4] data[n]
type Conform struct {
	lenlen int
	len    int
	Item
}

type ASCII struct {
	v    []byte
	null bool
}

type UTF16 struct {
	v    []byte
	null bool
}

type PNG struct{
	v []byte
}

func (t Repeat) Bytes()  []byte{ return bytes.Repeat([]byte{t.b}, t.n) }
func (t Run)    Bytes()  []byte{
	if len(t.items) == 0{
		return t.v
	}
	buf := new(bytes.Buffer)
	for _, v := range t.items{
		buf.Write(v.Bytes())
	}
	return buf.Bytes()
}
func (t ASCII)  Bytes()  []byte{ return t.v }
func (t UTF16)  Bytes()  []byte{ return t.v }
func (t PNG)    Bytes()  []byte{ return t.v }
func (t EOF)    Bytes()  []byte{ return nil }
func (t Conform) Bytes()   []byte{
	x := make([]byte, t.len+t.lenlen)
	copy(x[t.lenlen:], t.Item.Bytes())
	switch t.lenlen{
	case 4:bo.P32l(x,int32(t.len))
	case 2:bo.P16l(x,int16(t.len))
	case 1:x[0] = byte(t.len)
	}
 	return x
}
func (t Number) Bytes() []byte{
	x := make([]byte, t.width)
	switch t.width{
	case 4: bo.P32l(x, int32(t.v))
	case 3: panic("bug: t.Number.Bytes()")
	case 2: bo.P16l(x, int16(t.v))
	case 1: x[0] = byte(t.v)
	}
    return x
}

// LastN returns the last N bytes
func (a *ASCII) LastN(n int) []byte{
	l := len(a.v)
	if l <= n{
		return nil
	}
	return a.v[l-n:]
}

// LastN returns the last N bytes
func (a *UTF16) LastN(n int) []byte{
	l := len(a.v)
	if l <= n{
		return nil
	}
	return a.v[l-n:]
}

func (a *ASCII) String() string{
	return fmt.Sprintf("%#v", a)
}
func (u *UTF16) String() string {
	return fmt.Sprintf("%#v", u)
}
