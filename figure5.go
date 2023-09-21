// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// figure5 illustrates figure 5 of Featherweight Go.
//
// Unfortunately, it does not compile because the interface Eq
// constrains its own type parameter.
package main

import "fmt"

// Featherweight Go, Fig. 5

// Eq is the Eq interface from figure 5.
type Eq[a Eq[a]] interface {
	Equal(that a) bool
}

type Int int

func (this Int) Equal(that Int) bool {
	return this == that
}

// Pair is the Pair struct from figure 5, but with the type constraints
// strengthened to Eq because the final Go design for type parameters
// does not allow methods that refine the receiver's type constraints.
type Pair[a Eq, b Eq] struct {
	left  a
	right b
}

func (this Pair[a, b]) Equal(that Pair[a, b]) bool {
	return this.left.Equal(that.left) && this.right.Equal(that.right)
}

func main() {
	var i, j Int = 1, 2
	var p Pair[Int, Int] = Pair[Int, Int]{i, j}
	fmt.Println(p.Equal(p)) // true
}
