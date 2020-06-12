package main

import (
	"context"
	"fmt"
)

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

	// cannot use cancelSlice (type []context.CancelFunc) as type []func() in append
	// fc2 := append(funcSlice, cancelSlice...)
	// fmt.Printf("append(%T, %T...) = %T\n", funcSlice, cancelSlice, fc2)

	// cannot use cancels (type Cancels) as type []func() in append
	// FC2 := append(funcs, cancels...)
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

	// cannot use bidiSlice (type []chan int) as type []<-chan int in append
	// rb2 := append(recvSlice, bidiSlice...)
	// fmt.Printf("append(%T, %T...) = %T\n", recvSlice, bidiSlice, rb2)
}
