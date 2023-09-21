// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// atomic_value illustrates a type whose only requirement on its type parameter
// is that the argument be a concrete (non-interface) type.
//
// This requirement cannot be expressed in the current design, leading to
// the potential for run-time panics.
package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
)

// Value is a typed version of atomic.Value.
//
// But you'd better not instantiate it with an interface type, or you'll have to
// be very careful to only ever store values of a consistent concrete type!
type Value[T any] struct {
	a atomic.Value
}

func (v *Value[T]) Load() (x T) {
	return v.a.Load().(T)
}

func (v *Value[T]) Store(x T) {
	v.a.Store(x)
}

func main() {
	var err Value[error]
	err.Store(os.ErrNotExist)
	err.Store(syscall.ENOSYS) // Generics provide type-safety, right? 😉
	fmt.Printf("stored error: %v\n", err.Load())
}
