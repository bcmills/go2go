// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package iterate provides iterator adaptors for built-in container types.
package iterate

import (
	"context"
	"errors"
	"io"
	"reflect"
)

// Iterate invokes emit for each value produced by next.
// If next returns io.EOF, Iterate returns nil.
// If next returns any other non-nil error, Iterate immediately returns that error.
func Iterate[T any](next func() (T, error), emit func(T) error) error {
	for {
		x, err := next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			return nil
		}
		if err := emit(x); err != nil {
			return err
		}
	}
}

// Transform transforms the values produced by next by applying f to each value.
func Transform[T, U any](next func() (T, error), f func(T) (U, error)) func() (U, error) {
	return func () (U, error) {
		x, err := next()
		if err != nil {
			return *new(U), err
		}
		return f(x)
	}
}

// Reduce iterates next until EOF, accumulating the results using f with initial value init.
func Reduce[T, U any](next func() (T, error), init U, f func(T, U) (U, error)) (U, error) {
	y := init
	for {
		x, err := next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return y, err
			}
			break
		}
		y, err = f(x, y)
		if err != nil {
			return y, err
		}
	}
	return y, nil
}

// SliceElems returns an iterator function that produces each successive element of s.
func SliceElems[T any](s []T) func() (T, error) {
	i := 0
	return func() (T, error) {
		if i >= len(s) {
			return *new(T), io.EOF
		}

		x := s[i]
		i++
		return x, nil
	}
}

// ToSlice returns a slice containing the elements produced by next.
func ToSlice[T any](next func() (T, error)) ([]T, error) {
	var s []T
	for {
		x, err := next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return s, err
			}
			return s, nil
		}
		s = append(s, x)
	}
}

// MapKeys returns an iterator function that produces each successive element
// of m, in arbitrary order.
func MapKeys[K comparable, V any](m map[K]V) func() (K, error) {
	iter := reflect.ValueOf(m).MapRange()
	if !iter.Next() {
		iter = nil
	}
	return func() (K, error) {
		if iter == nil {
			return *new(K), io.EOF
		}
		k, _ := iter.Key().Interface().(K)
		if !iter.Next() {
			iter = nil
		}
		return k, nil
	}
}

// MapElems returns an iterator function that produces each successive element
// of m, in arbitrary order.
func MapElems[K comparable, V any](m map[K]V) func() (V, error) {
	iter := reflect.ValueOf(m).MapRange()
	if !iter.Next() {
		iter = nil
	}
	return func() (V, error) {
		if iter == nil {
			return *new(V), io.EOF
		}
		v, _ := iter.Value().Interface().(V)
		if !iter.Next() {
			iter = nil
		}
		return v, nil
	}
}

// Map returns an iterator function that produces each successive key–value pair
// of m, in arbitrary order.
func Map[K comparable, V any](m map[K]V) func() (K, V, error) {
	iter := reflect.ValueOf(m).MapRange()
	if !iter.Next() {
		iter = nil
	}
	return func() (K, V, error) {
		if iter == nil {
			return *new(K), *new(V), io.EOF
		}
		k, _ := iter.Key().Interface().(K)
		v, _ := iter.Value().Interface().(V)
		if !iter.Next() {
			iter = nil
		}
		return k, v, nil
	}
}

// Chan returns an iterator function that returns each successive element
// received from c.
func Chan[T any](c <-chan T) func() (T, error) {
	return func() (T, error) {
		x, ok := <- c
		if !ok {
			return x, io.EOF
		}
		return x, nil
	}
}

// Chan returns an iterator function that returns each successive element
// received from c, or the zero T and ctx.Err() if ctx is done and no value is
// ready to receive.
func ChanCtx[T any](ctx context.Context, c <-chan T) func() (T, error) {
	return func() (T, error) {
		select {
		case <-ctx.Done():
			return *new(T), ctx.Err()
		case x, ok := <-c:
			if !ok {
				return x, io.EOF
			}
			return x, nil
		}
	}
}

// ToChan repeatedly produces a value using next and sends that value on c,
// until next returns a non-nil error.
//
// ToChan does not close the channel.
// The caller may close it after ToChan returns.
//
// If the error returned by next is io.EOF, ToChan returns nil.
// Otherwise, ToChan returns the error returned by next.
func ToChan[T any](c chan<- T, next func() (T, error)) error {
	for {
		x, err := next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			return nil
		}

		c <- x
	}
}

// ToChanCtx repeatedly produces a value using next and sends that value on c,
// until either ctx is done or next returns a non-nil error.
//
// ToChanCtx does not close the channel.
// The caller may close it after ToChanCtx returns.
//
// If the final error from next is io.EOF,
// ToChanCtx returns the final value from next and a nil error.
// If ctx becomes done after a call to next returned a nil error,
// ToChanCtx returns the final (unsent) value from next and ctx.Err().
// Otherwise, ToChanCtx returns the final value and error from next.
//
func ToChanCtx[T any](ctx context.Context, c chan<- T, next func() (T, error)) (T, error) {
	for {
		x, err := next()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return x, err
			}
			return x, nil
		}

		select {
		case <-ctx.Done():
			return x, ctx.Err()
		case c <- x:
		}
	}
}
