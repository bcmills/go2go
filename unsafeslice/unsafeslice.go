// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package unsafeslice implements common unsafe transformations involving slices,
// based on the reflective prototype in github.com/bcmills/unsafeslice.
//
// unsafeslice uses the reflect package only for the SliceHeader type definition.
// That dependency could be eliminated using a local definition and regression
// test, as is done in the internal/unsafeheader package in Go 1.15.
package unsafeslice

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Of returns a slice of length and capacity n located at p.
//
// The caller must ensure that p points to a backing array containing at least n
// elements with an equivalent layout, size, and alignment to T.
//
// This implements one possible API for https://golang.org/issue/19367
// and https://golang.org/issue/13656.
func Of[T any](p *T, n int) []T {
	return unsafe.Slice(p, n)
}

// Convert returns a slice that refers to the same memory region as the slice src,
// but at an arbitrary element type.
//
// At some call sites, ConvertAt may provide better type inference than Convert.
//
// The caller must ensure that the length and capacity in bytes of src are
// integer multiples of the size of T2, and that the fields at each byte offset
// in the resulting slices have equivalent layouts.
//
// This implements one possible API for https://golang.org/issue/38203.
func Convert[T1, T2 any](src []T1) []T2 {
	srcElemSize := reflect.TypeOf(src).Elem().Size()
	capBytes := uintptr(cap(src)) * srcElemSize
	lenBytes := uintptr(len(src)) * srcElemSize

	dstElemSize := reflect.TypeOf((*T2)(nil)).Elem().Size()

	if capBytes%dstElemSize != 0 {
		panic(fmt.Sprintf("Convert: src capacity (%d bytes) is not a multiple of dst element size (%T: %d bytes)", capBytes, *new(T2), dstElemSize))
	}
	dstCap := capBytes / dstElemSize
	if int(dstCap) < 0 || uintptr(int(dstCap)) != dstCap {
		panic(fmt.Sprintf("Convert: dst capacity (%d) overflows int", dstCap))
	}

	if lenBytes%dstElemSize != 0 {
		panic(fmt.Sprintf("Convert: src length (%d bytes) is not a multiple of dst element size (%T: %d bytes)", lenBytes, *new(T2), dstElemSize))
	}
	dstLen := lenBytes / dstElemSize
	if int(dstLen) < 0 || uintptr(int(dstLen)) != dstLen {
		panic(fmt.Sprintf("ConvertAt: dst length (%d) overflows int", dstLen))
	}

	return unsafe.Slice((*T2)(unsafe.Pointer(unsafe.SliceData(src))), dstCap)[:dstLen]
}

// ConvertAt sets dst, which must be non-nil, to a slice that refers to the same
// memory region as the slice src, but possibly at a different type.
//
// The caller must ensure that the length and capacity in bytes of src are
// integer multiples of the size of T2, and that the fields at each byte offset
// in the resulting slices have equivalent layouts.
//
// This implements one possible API for https://golang.org/issue/38203.
func ConvertAt[T2, T1 any](dst *[]T2, src []T1) {
	*dst = Convert[T1, T2](src)
}

// AsPointer returns a pointer to the array backing src[0:len(src)] as type *T.
//
// The caller must ensure that the length in bytes of src is an integer multiple
// of the size of T, and that the fields at each byte offset in the resulting
// slice have a layout equivalent to T.
//
// At some call sites, SetPointer may provide better type inference than
// AsPointer.
func AsPointer[E, T any](src []E) *T {
	return unsafe.SliceData(Convert[E, T](src[:len(src):len(src)]))
}

// SetPointer sets dst, which must be non-nil, to a pointer that refers to the
// elements of src. Typically, dst should point to a pointer to an array with
// the same length and element type as src.
func SetPointer[T, E any](dst **T, src []E) {
	*dst = AsPointer[E, T](src)
}
