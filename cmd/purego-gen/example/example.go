// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Ebitengine Authors

//go:build (darwin || freebsd || linux || netbsd) && (amd64 || arm64)

// Package example demonstrates purego-gen usage.
//
// Usage:
//
//	//go:generate go run github.com/ebitengine/purego/cmd/purego-gen
//
// Then define your functions with //purego:sym annotations:
//
//	//purego:sym lib symbol_name
//	var myFunc func(int32, unsafe.Pointer, int32) int32
//
// Call RegisterPuregoFuncs() after loading your library to set up
// all the zero-allocation function wrappers.
package example

import "unsafe"

// lib is the loaded library handle - set this before calling RegisterPuregoFuncs()
var lib uintptr

// Simple function with integer args and return
//
//purego:sym lib bench_add
var benchAdd func(int64, int64) int64

// Function with pointer argument (KeepAlive will be generated)
//
//purego:sym lib bench_ptr
var benchPtr func(unsafe.Pointer, int64) int64

// Function with bool arg and return
//
//purego:sym lib is_valid
var isValid func(int32) bool

// Function with multiple pointer types
//
//purego:sym lib opus_encode
var opusEncode func(uint64, unsafe.Pointer, int32, unsafe.Pointer, int32) int32

// Function with typed pointers
//
//purego:sym lib process_data
var processData func(*int32, *byte, int64) int32

// Void function (no return)
//
//purego:sym lib no_return_fn
var noReturnFn func(int64, int64, int64)

// Function with many arguments (uses Syscall9)
//
//purego:sym lib many_args
var manyArgs func(int64, int64, int64, int64, int64, int64, int64, int64, int64) int64

// Function with bool argument
//
//purego:sym lib bool_and
var boolAnd func(bool, bool) bool

//go:generate go run github.com/ebitengine/purego/cmd/purego-gen
