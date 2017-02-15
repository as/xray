package main
import "fmt"
type Item interface{}

type EOF struct {
	err error
}
type Repeat struct {
	b byte
	n int
}

type Run struct {
	s     []byte
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
	return fmt.Sprintf("ascii=%q", a.v)
}
func (u *UTF16) String() string {
	return fmt.Sprintf("utf16=%q", u.v)
}
