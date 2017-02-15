package main

import (
	"fmt"
)

func (p *parser) printTrace(a ...interface{}) {
	return
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

type tab int

func (t tab) pre() {
	for i := 0; i < int(t); i++ {
		fmt.Print("\t")
	}
}
func (t tab) Printf(fm string, i ...interface{}) { t.pre(); fmt.Printf(fm, i...) }
func (t tab) Println(i ...interface{})           { t.pre(); fmt.Println(i...) }

var Tab tab

func println(i Item) {
	print(i)
	fmt.Println()
	fmt.Println()
}

func print(i Item) {
	switch t := i.(type) {
	case *Repeat:
		Tab.Printf("%#v\n", t)
	case *Run:
		//Tab.Printf("Run: bytes=%q\n", t.s)
		Tab++
		for _, v := range t.items {
			print(v)
		}
		Tab--
	case *Conform:
		Tab.Printf("%#v\n", t)
		Tab++
		print(t.Item)
		Tab--
	case fmt.Stringer:
		Tab.Printf("%s\n", t.String())
	case Item:
		Tab.Printf("%#v \n", t)
	}
}

func trace(p *parser, msg string) *parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."))
func un(p *parser) {
	p.indent--
	p.printTrace(")")
}
