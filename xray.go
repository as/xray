// Copyright 2017 "as". All rights reserved. Torgo is governed
// the same BSD license as the go programming language.

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/as/bo"
)

func Gintl(b []byte) int64 {
	return int64(bo.Gintl(b))
}

// conformsAny checks if n represents the size of the Item
func conformsAny(it Item, n int) (lenlen, width int){
	mask := 0xffffffff<<8
	lenlen = 1
	for mask != 0{
		if c := conformsTo(it, n&^mask); c != 0{
			return lenlen, n&^mask
		}
		lenlen++
		mask <<= 8
	}
	fmt.Println("lenlen=%x n=%x mask=%x n&^mask=%x\n", lenlen, n, mask, n&^mask)
	fmt.Println("*********************")
	return 0, 0
}

func conformsTo(i Item, n int) (width int) {
	switch t := i.(type) {
	case *ASCII:
		if n == len(t.v) {
			return 1
		}	
		if n == len(t.v)-1 && zero(t.LastN(1)) {
			// NUL not included in conforming length
			return 1
		}							
	case *UTF16:
		if n == len(t.v)*2 {
			return 1
		}	
		if n == (len(t.v)-1)*2 && zero(t.LastN(2)) {
			return 1
		}
	}
	return 0
}

func zero(b []byte) bool{
	if len(b) == 0{
		return false
	}
	for _, v := range b{
		if v != 0{
			return false
		}
	}
	return true
}

const Common = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()-_+={[]}\;',./<>?`

func common(b byte) bool {
	for _, v := range Common {
		if b == byte(v) {
			return true
		}
	}
	return false
}

func main() {
	p := &parser{br: bufio.NewReader(os.Stdin)}
	items := p.parse()
	for i, v := range items {
		fmt.Printf("Item %d: ", i)
		print(v)
		switch t := v.(type) {
		case *Run:
			i := 0
			for i = range t.s {
				if common(t.s[i]) {
					break
				}
			}
			if i == len(t.s) {
				// try integer types
			} else {
				// may be a string
			}
		}
	}

}
