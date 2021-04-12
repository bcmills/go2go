// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package containers defines generic interfaces for built-in container types.
package containers

type Lenner interface {
	Len() int
}

type Capper interface {
	Cap() int
}

type Indexer[K, V any] interface {
	Index(K) (V, bool)
}

type IndexSetter[K, V any] interface {
	Index(K) (V, bool)
	SetIndex(K, V)
}

type Sender[V any] interface {
	Send(V)
}

type Closer interface {
	Close()
}

type Receiver[V any] interface {
	Receive() (V, bool)
}

type KeyRanger[K any] interface {
	RangeKeys(func(K) (ok bool))
}

type ElemRanger[V any] interface {
	RangeElems(func(V) (ok bool))
}

type Ranger[K, V any] interface {
	RangeKeys(func(K) (ok bool))
	RangeElems(func(V) (ok bool))
	Range(func(K, V) (ok bool))
}

type String string

func (s String) Len() int { return len(s) }

func (s String) Index(i int) (byte, bool) {
	if i < 0 || i >= len(s) {
		return 0, false
	}
	return s[i], true
}

func (s String) RangeKeys(f func(i int) bool) {
	for i := range s {
		if !f(i) {
			break
		}
	}
}

func (s String) RangeElems(f func(r rune) bool) {
	for _, r := range s {
		if !f(r) {
			break
		}
	}
}

func (s String) Range(f func(i int, r rune) bool) {
	for i, r := range s {
		if !f(i, r) {
			break
		}
	}
}

type Slice[T any] []T

func (s Slice[T]) Len() int { return len(s) }
func (s Slice[T]) Cap() int { return cap(s) }
func (s Slice[T]) SetIndex(i int, x T) { s[i] = x }

func (s Slice[T]) Index(i int) (T, bool) {
	if i < 0 || i >= len(s) {
		return *new(T), false
	}
	return s[i], true
}

func (s Slice[T]) RangeKeys(f func(i int) bool) {
	for i := range s {
		if !f(i) {
			break
		}
	}
}

func (s Slice[T]) RangeElems(f func(x T) bool) {
	for _, x := range s {
		if !f(x) {
			break
		}
	}
}

func (s Slice[T]) Range(f func(i int, x T) bool) {
	for i, x := range s {
		if !f(i, x) {
			break
		}
	}
}

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Len() int { return len(m) }
func (m Map[K, V]) SetIndex(k K, v V) { m[k] = v }

func (m Map[K, V]) Index(k K) (V, bool) {
	v, ok := m[k]
	return v, ok
}

func (m Map[K, V]) RangeKeys(f func(K) bool) {
	for k := range m {
		if !f(k) {
			break
		}
	}
}

func (m Map[K, V]) RangeElems(f func(V) bool) {
	for _, v := range m {
		if !f(v) {
			break
		}
	}
}

func (m Map[K, V]) Range(f func(K, V) bool) {
	for k, v := range m {
		if !f(k, v) {
			break
		}
	}
}



type Chan[T any] chan T

func (c Chan[T]) Len() int { return len(c) }
func (c Chan[T]) Cap() int { return cap(c) }
func (c Chan[T]) Send(x T) { c <- x }
func (c Chan[T]) Close() { close(c) }

func (c Chan[T]) Recv() (T, bool) {
	x, ok := <-c
	return x, ok
}

func (c Chan[T]) RangeElems(f func(T) bool) {
	for x := range c {
		if !f(x) {
			break
		}
	}
}

type RecvChan[T any] <-chan T

func (c RecvChan[T]) Len() int { return len(c) }
func (c RecvChan[T]) Cap() int { return cap(c) }

func (c RecvChan[T]) Recv() (T, bool) {
	x, ok := <-c
	return x, ok
}

func (c RecvChan[T]) RangeElems(f func(x T) bool) {
	for x := range c {
		if !f(x) {
			break
		}
	}
}

type SendChan[T any] chan<- T

func (c SendChan[T]) Len() int { return len(c) }
func (c SendChan[T]) Cap() int { return cap(c) }
func (c SendChan[T]) Send(x T) { c <- x }
func (c SendChan[T]) Close() { close(c) }
