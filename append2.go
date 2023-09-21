// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// append2 illustrates the properties of a two-type-parameter "append" variant.
// With explicit type annotations, it seems to be able to handle the same cases
// as the existing built-in append.
//
// However, I don't see a plausible type-inference algorithm that makes it work
// for all of those cases without explicit type annotations.
package main

import (
	"context"
	"fmt"
)

type sliceOf[E any] interface{ ~[]E }

func append[T any, S sliceOf[T]](s S, t ...T) S {
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
	r           <-chan int
	recvSlice   []<-chan int
	R           Recv
	RecvSlice   []Recv
	b           chan int
	bidiSlice   []chan int
)

func main() {
	ff := append(funcSlice, f)
	fmt.Printf("append(%T, %T) = %T\n", funcSlice, f, ff)

	Ff := append(funcs, f)
	fmt.Printf("append(%T, %T) = %T\n", funcs, f, Ff)

	// []func() does not satisfy sliceOf[context.CancelFunc] ([]func() missing in ~[]context.CancelFunc)
	fc := append[func()](funcSlice, cancel)
	fmt.Printf("append(%T, %T) = %T\n", funcSlice, cancel, fc)

	cf := append(cancelSlice, f)
	fmt.Printf("append(%T, %T) = %T\n", cancelSlice, f, cf)

	// Funcs does not satisfy sliceOf[context.CancelFunc] (Funcs missing in ~[]context.CancelFunc)
	Fc := append[func()](funcs, cancel)
	fmt.Printf("append(%T, %T) = %T\n", funcs, cancel, Fc)

	Cc := append(cancels, f)
	fmt.Printf("append(%T, %T) = %T\n", cancels, f, Cc)

	// []func() does not satisfy sliceOf[context.CancelFunc] ([]func() missing in ~[]context.CancelFunc)
	ffc := append[func()](funcSlice, f, cancel)
	fmt.Printf("append(%T, %T, %T) = %T\n", funcSlice, f, cancel, ffc)

	ff2 := append(funcSlice, funcSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", funcSlice, funcSlice, ff2)

	FF2 := append(funcs, funcs...)
	fmt.Printf("append(%T, %T...) = %T\n", funcs, funcs, FF2)

	Ff2 := append(funcs, funcSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", funcs, funcSlice, Ff2)

	// []func() does not satisfy sliceOf[context.CancelFunc] ([]func() missing in ~[]context.CancelFunc)
	// cannot use cancelSlice (variable of type []context.CancelFunc) as []func() value in argument to append[func()]
	// fc2 := append[func()](funcSlice, cancelSlice...)
	// fmt.Printf("append(%T, %T...) = %T\n", funcSlice, cancelSlice, fc2)

	// Funcs does not satisfy sliceOf[context.CancelFunc] (Funcs missing in ~[]context.CancelFunc)
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

	// cannot use bidiSlice (variable of type []chan int) as []<-chan int value in argument to append
	// rb2 := append(recvSlice, bidiSlice...)
	// fmt.Printf("append(%T, %T...) = %T\n", recvSlice, bidiSlice, rb2)
}
