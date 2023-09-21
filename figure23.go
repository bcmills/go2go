// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

// Adapted from Featherweight Go, Fig. 23

type Eq[a any] interface {
	Equal(that a) Bool
}
type Bool interface {
	Not() Bool
	Equal(that Bool) Bool
	cond(br Branches[any]) any
}
type Branches[a any] interface {
	IfTT() a
	IfFF() a
}
type Cond[a any] struct{}
func (this Cond[a]) Eval(b Bool, br Branches[a]) a {
	return b.cond(anyBranches[a]{br}).(a)
}
type anyBranches[a any] struct {
	typed Branches[a]
}
func (br anyBranches[a]) IfTT() any {
	return br.typed.IfTT()
}
func (br anyBranches[a]) IfFF() any {
	return br.typed.IfFF()
}

type TT struct{}
type FF struct{}

func (this TT) Not() Bool { return FF{} }
func (this FF) Not() Bool { return TT{} }

func (this TT) Equal(that Bool) Bool { return that }
func (this FF) Equal(that Bool) Bool { return that.Not() }

func (this TT) cond(br Branches[any]) any { return br.IfTT() }
func (this FF) cond(br Branches[any]) any { return br.IfFF() }

func main() {}
