// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// figure5 illustrates figure 5 of Featherweight Go,
// with the generic types defined in figure 5 converted to plain interfaces.
package main

// Featherweight Go, Fig. 3

type Function[a any, b any] interface {
	Apply(x a) b
}
type incr struct{ n int }
func (this incr) Apply(x int) int {
	return x + this.n
}
type pos struct{}
func (this pos) Apply(x int) bool {
	return x > 0
}
type compose[a any, b any, c any] struct {
	f Function[a, b]
	g Function[b, c]
}
func (this compose[a, b, c]) Apply(x a) c {
	return this.g.Apply(this.f.Apply(x))
}

// Adapted from Featherweight Go, Fig. 4

type Eq[a any] interface {
 Equal(a) bool
}
type Int int
func (this Int) Equal(that Int) bool {
	return this == that
}
type List[a any] interface{
	Match(casenil Function[Nil[a], any], casecons Function[Cons[a], any]) any
}
type Nil[a any] struct{}
func (xs Nil[a]) Match(casenil Function[Nil[a], any], casecons Function[Cons[a], any]) any {
	return casenil.Apply(xs)
}
type Cons[a any] struct {
	Head a
	Tail List[a]
}
func (xs Cons[a]) Match(casenil Function[Nil[a], any], casecons Function[Cons[a], any]) any {
	return casecons.Apply(xs)
}
type lists[a any, b any] struct{}
func (_ lists[a, b]) Map(f Function[a, b], xs List[a]) List[b] {
	return xs.Match(mapNil[a, b]{}, mapCons[a, b]{f}).(List[b])
}
type mapNil[a any, b any] struct{}
func (m mapNil[a, b]) Apply(_ Nil[a]) any {
	return Nil[b]{}
}
type mapCons[a any, b any] struct{
	f Function[a, b]
}
func (m mapCons[a, b]) Apply(xs Cons[a]) any {
	return Cons[b]{m.f.Apply(xs.Head), lists[a, b]{}.Map(m.f, xs.Tail)}
}

// Featherweight Go, Fig. 5 (monomorphized)

type Edge interface {
	Source() Vertex
	Target() Vertex
}
type Vertex interface {
	Edges() List[Edge]
}

func main() {}
