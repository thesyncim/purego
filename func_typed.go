// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build (darwin || freebsd || linux || netbsd) && (amd64 || arm64)

//go:generate go run ./internal/codegen/gen_typed.go

package purego

import (
	"math"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego/internal/strings"
)

// tryRegisterTyped attempts to register a function with zero per-call allocations.
// It uses type assertions to match common function signatures and create specialized
// call paths that bypass reflect.MakeFunc entirely.
//
// Supported signatures (zero alloc):
//   - func() R
//   - func(A) R
//   - func(A, B) R
//   - func(A, B, C) R
//   - ... up to 8 args
//
// Where A, B, C, R are: int, int32, int64, uint, uint32, uint64, uintptr,
// unsafe.Pointer, float32, float64
//
// Returns true if an optimized path was used, false otherwise.

func tryRegisterTyped(fptr any, cfn uintptr) bool {
	// Try generated integer/pointer/bool signatures first (most common)
	if tryRegisterTypedGenerated(fptr, cfn) {
		return true
	}

	// Handle special cases: floats (need float registers) and strings (need allocation)
	switch fn := fptr.(type) {
	// Float returns (0 args)
	case *func() float32:
		*fn = func() float32 { return syscall0Float(cfn) }
		return true
	case *func() float64:
		*fn = func() float64 { return syscall0Double(cfn) }
		return true

	// Float args (1 arg)
	case *func(float64) float64:
		*fn = func(a float64) float64 { return syscall1Double(cfn, a) }
		return true
	case *func(float32) float32:
		*fn = func(a float32) float32 { return syscall1Float(cfn, a) }
		return true

	// String args (1 alloc for non-null-terminated strings)
	// Must use runtime.KeepAlive to prevent GC from collecting the C string
	// before the syscall completes
	case *func(string) int:
		*fn = func(s string) int {
			cstr := strings.CString(s)
			r := int(Syscall1(cfn, uintptr(unsafe.Pointer(cstr))))
			runtime.KeepAlive(cstr)
			return r
		}
		return true
	case *func(string) uintptr:
		*fn = func(s string) uintptr {
			cstr := strings.CString(s)
			r := Syscall1(cfn, uintptr(unsafe.Pointer(cstr)))
			runtime.KeepAlive(cstr)
			return r
		}
		return true
	case *func(string):
		*fn = func(s string) {
			cstr := strings.CString(s)
			Syscall1(cfn, uintptr(unsafe.Pointer(cstr)))
			runtime.KeepAlive(cstr)
		}
		return true
	}

	return false
}

// Float syscall helpers

func syscall0Float(fn uintptr) float32 {
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := math.Float32frombits(uint32(syscall.f1))
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

func syscall0Double(fn uintptr) float64 {
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := math.Float64frombits(uint64(syscall.f1))
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

func syscall1Float(fn uintptr, a float32) float32 {
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.f1 = uintptr(math.Float32bits(a))
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := math.Float32frombits(uint32(syscall.f1))
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

func syscall1Double(fn uintptr, a float64) float64 {
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.f1 = uintptr(math.Float64bits(a))
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := math.Float64frombits(uint64(syscall.f1))
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}
