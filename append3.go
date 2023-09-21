// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// append3 illustrates the properties of a 3-type-parameter "append" variant.
//
// With a more advanced type-inference algorithm and a proper "assignable to"
// constraint, it could support inference for all of the same cases as the
// built-in "append" does today, plus a few others.
package main

import (
	"context"
	"fmt"
	"reflect"
)

type sliceOf[E any] interface{ ~[]E }

// buggyAssignableTo simulates an “assignable to” type constraint using
// something a little too permissive ("any").
// We confirm assignability at run-time using the reflect package.
type buggyAssignableTo[T any] interface { any }

func append[T any, T2 buggyAssignableTo[T], S sliceOf[T]](s S, t ...T2) S {
	// Confirm that T2 is assignable to T.
	// Ideally this should happen in the type system instead of at run-time.
	rt := reflect.TypeOf(s).Elem()
	rt2 := reflect.TypeOf(t).Elem()
	if !rt2.AssignableTo(rt) {
		panic("append: T2 is not assignable to T")
	}

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
	for i, x := range t {
		// We need to bounce through reflect because buggyAssignableTo doesn't
		// actually enable assignment.
		xt := reflect.ValueOf(x).Convert(rt).Interface().(T)
		s[lens+i] = xt
	}
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

	fc := append(funcSlice, cancel)
	fmt.Printf("append(%T, %T) = %T\n", funcSlice, cancel, fc)

	cf := append(cancelSlice, f)
	fmt.Printf("append(%T, %T) = %T\n", cancelSlice, f, cf)

	Fc := append(funcs, cancel)
	fmt.Printf("append(%T, %T) = %T\n", funcs, cancel, Fc)

	Cc := append(cancels, f)
	fmt.Printf("append(%T, %T) = %T\n", cancels, f, Cc)

	ffc := append(funcSlice, f, cancel)
	fmt.Printf("append(%T, %T, %T) = %T\n", funcSlice, f, cancel, ffc)

	ff2 := append(funcSlice, funcSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", funcSlice, funcSlice, ff2)

	FF2 := append(funcs, funcs...)
	fmt.Printf("append(%T, %T...) = %T\n", funcs, funcs, FF2)

	Ff2 := append(funcs, funcSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", funcs, funcSlice, Ff2)

	fc2 := append(funcSlice, cancelSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", funcSlice, cancelSlice, fc2)

	FC2 := append(funcs, cancels...)
	fmt.Printf("append(%T, %T...) = %T\n", funcs, cancels, FC2)

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

	rb2 := append(recvSlice, bidiSlice...)
	fmt.Printf("append(%T, %T...) = %T\n", recvSlice, bidiSlice, rb2)
}
