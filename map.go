// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// map illustrates the awkwardness of functional APIs in the absence of
// assignability and/or subtyping constraints.
package main

// Map returns a new slice containing the result of applying f to each element
// in src.
func Map[T1, T2 any](src []T1, f func(T1) T2) []T2 {
	dst := make([]T2, 0, len(src))
	return AppendMapped(dst, src, f)
}

// AppendMapped applies f to each element in src, appending each result to dst
// and returning the resulting slice.
func AppendMapped[T1, T2 any](dst []T2, src []T1, f func(T1) T2) []T2 {
	for _, x := range src {
		dst = append(dst, f(x))
	}
	return dst
}

func recv[T any](x <-chan T) T {
	return <-x
}

func main() {
	var chans []chan int

	// To map the recv function over a slice of bidirectional channels,
	// we need to wrap it: even though the element type "chan int"
	// is assignable to the argument type, the two function types are distinct.
	//
	// (There is no way for Map to declare an argument type “assignable from T1”,
	// so it must use “exactly T1” instead.)
	vals := Map(chans, func(x chan int) int {
		return recv(x)
	})

	_ = vals
}
