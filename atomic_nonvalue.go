// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// atomic_nonvalue illustrates a (less efficient) workaround for the atomic_value
// example.
//
// It eliminates the type-safety hole by inducing an extra allocation that is
// normally not needed.
package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
)

type Value[T any] struct {
	a atomic.Pointer[T]
}

func (v *Value[T]) Store(x T) {
	// Store a pointer to x instead of x itself, in case T is an interface type.
	//
	// It would be more efficient to require the caller to instantiate Value with
	// a non-interface type, so that the extra pointer allocation would be visible
	// on the caller side and could be avoided for non-interface T.
	v.a.Store(&x)
}

func (v *Value[T]) Load() (x T) {
	return *v.a.Load()
}

func main() {
	var err Value[error]
	err.Store(os.ErrNotExist)
	err.Store(syscall.ENOSYS)
	fmt.Printf("stored error: %v\n", err.Load())
}
