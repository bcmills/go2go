// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// append1 illustrates the properties of the "append" variation described in the
// Type Parameters draft design.
package main

import (
	"context"
	"fmt"
)

func append[T any](s []T, t ...T) []T {
	lens := len(s)
	tot := lens + len(t)
	if tot < 0 {
		panic("append: cap out of range")
	}
	if tot > cap(s) {
		news := make([]T, tot, tot+tot/2)
		copy(news, s)
		s = news
	}
	s = s[:tot]
	copy(s[lens:], t)
	return s
}

type Funcs []func()
type Cancels []context.CancelFunc
type Recv <-chan int

var (
	f           func()
	cancel      context.CancelFunc
	funcSlice   []func()
	cancelSlice []context.CancelFunc
	funcs       Funcs
	cancels     Cancels
	r        <-chan int
	recvSlice   []<-chan int
	R        Recv
	RecvSlice   []Recv
	b        chan int
	bidiSlice   []chan int
)

func main() {
	ff := append(funcSlice, f)
	fmt.Printf("append(%T, %T) = %T\n", funcSlice, f, ff)

	// returns type []func() instead of Funcs
	Ff := append(funcs, f)
	fmt.Printf("append(%T, %T) = %T\n", funcs, f, Ff)

	// cannot use funcSlice (variable of type []func()) as []context.CancelFunc value in argument to append
	fc := append[func()](funcSlice, cancel)
	fmt.Printf("append(%T, %T) = %T\n", funcSlice, cancel, fc)

	cf := append(cancelSlice, f)
	fmt.Printf("append(%T, %T) = %T\n", cancelSlice, f, cf)

	// returns type []func instead of Funcs
	// cannot use funcs (variable of type Funcs) as []context.CancelFunc value in argument to append
	Fc := append[func()](funcs, cancel)
	fmt.Printf("append(%T, %T) = %T\n", funcs, cancel, Fc)

	Cc := append(cancels, f)
	fmt.Printf("append(%T, %T) = %T\n", cancels, f, Cc)

	// cannot use funcSlice (variable of type []func()) as []context.CancelFunc value in argument to append
	ffc := append[func()](funcSlice, f, cancel)
	fmt.Printf("append(%T, %T, %T) = %T\n", funcSlice, f, cancel, ffc)

	ff2 := append(funcSlice, funcSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", funcSlice, funcSlice, ff2)

	// returns type []func() instead of Funcs
	FF2 := append(funcs, funcs...)
	fmt.Printf("append(%T, %T...) = %T\n", funcs, funcs, FF2)

	// returns type []func() instead of Funcs
	Ff2 := append(funcs, funcSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", funcs, funcSlice, Ff2)

	// type []context.CancelFunc of cancelSlice does not match inferred type []func() for []T
	// cannot use cancelSlice (variable of type []context.CancelFunc) as []func() value in argument to append[func()]
	// fc2 := append[func()](funcSlice, cancelSlice...)
	// fmt.Printf("append(%T, %T...) = %T\n", funcSlice, cancelSlice, fc2)

	// type Cancels of cancels does not match inferred type []func() for []T
	// cannot use cancels (variable of type Cancels) as []func() value in argument to append[func()]
	// FC2 := append[func()](funcs, cancels...)
	// fmt.Printf("append(%T, %T...) = %T\n", funcs, cancels, FC2)

	rr := append(recvSlice, r)
	fmt.Printf("append(%T, %T) = %T\n", recvSlice, r, rr)

	rb := append(recvSlice, b)
	fmt.Printf("append(%T, %T) = %T\n", recvSlice, b, rb)

	RR := append(RecvSlice, R)
	fmt.Printf("append(%T, %T) = %T\n", RecvSlice, R, RR)

	Rb := append(RecvSlice, b)
	fmt.Printf("append(%T, %T) = %T\n", RecvSlice, b, Rb)

	rrb := append(recvSlice, r, b)
	fmt.Printf("append(%T, %T) = %T\n", recvSlice, b, rrb)

	rr2 := append(recvSlice, recvSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", recvSlice, recvSlice, rr2)

	// type []chan int of bidiSlice does not match inferred type []<-chan int for []T
	// cannot use bidiSlice (variable of type []chan int) as []<-chan int value in argument to append[<-chan int]
	// rb2 := append[<-chan int](recvSlice, bidiSlice...)
	// fmt.Printf("append(%T, %T...) = %T\n", recvSlice, bidiSlice, rb2)
}
