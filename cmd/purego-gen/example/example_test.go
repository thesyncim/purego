// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Ebitengine Authors

//go:build (darwin || freebsd || linux || netbsd) && (amd64 || arm64)

package example

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

var testLib uintptr

func TestMain(m *testing.M) {
	// Build C test library
	dir, err := os.MkdirTemp("", "purego-gen-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Create C source file
	cSrc := filepath.Join(dir, "test.c")
	err = os.WriteFile(cSrc, []byte(`
#include <stdint.h>
#include <stdbool.h>

int64_t bench_add(int64_t a, int64_t b) {
    return a + b;
}

int64_t bench_ptr(void* ptr, int64_t offset) {
    return (int64_t)ptr + offset;
}

bool is_valid(int32_t val) {
    return val > 0;
}

int32_t opus_encode(uint64_t encoder, void* pcm, int32_t frame_size, void* output, int32_t output_cap) {
    // Simulate opus encode - just return frame_size
    return frame_size;
}

int32_t process_data(int32_t* a, unsigned char* b, int64_t len) {
    return *a + (int32_t)*b + (int32_t)len;
}

void no_return_fn(int64_t a, int64_t b, int64_t c) {
    // Do nothing
}

int64_t many_args(int64_t a0, int64_t a1, int64_t a2, int64_t a3, int64_t a4, int64_t a5, int64_t a6, int64_t a7, int64_t a8) {
    return a0 + a1 + a2 + a3 + a4 + a5 + a6 + a7 + a8;
}

bool bool_and(bool a, bool b) {
    return a && b;
}
`), 0644)
	if err != nil {
		panic(err)
	}

	// Compile shared library
	libPath := filepath.Join(dir, "test.so")
	out, err := exec.Command("go", "env", "CC").Output()
	if err != nil {
		panic(err)
	}
	cc := string(out[:len(out)-1])
	if cc == "" {
		cc = "cc"
	}
	if err := exec.Command(cc, "-shared", "-o", libPath, cSrc).Run(); err != nil {
		panic(err)
	}

	// Load library
	testLib, err = load.OpenLibrary(libPath)
	if err != nil {
		panic(err)
	}
	lib = testLib

	// Register functions
	if err := RegisterPuregoFuncs(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestBenchAdd(t *testing.T) {
	result := benchAdd(10, 20)
	if result != 30 {
		t.Errorf("benchAdd(10, 20) = %d, want 30", result)
	}

	result = benchAdd(-5, 15)
	if result != 10 {
		t.Errorf("benchAdd(-5, 15) = %d, want 10", result)
	}
}

func TestBenchPtr(t *testing.T) {
	data := make([]byte, 100)
	ptr := unsafe.Pointer(&data[0])
	result := benchPtr(ptr, 50)
	expected := int64(uintptr(ptr)) + 50
	if result != expected {
		t.Errorf("benchPtr(%p, 50) = %d, want %d", ptr, result, expected)
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		input int32
		want  bool
	}{
		{1, true},
		{100, true},
		{0, false},
		{-1, false},
		{-100, false},
	}

	for _, tt := range tests {
		got := isValid(tt.input)
		if got != tt.want {
			t.Errorf("isValid(%d) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestOpusEncode(t *testing.T) {
	pcm := make([]int16, 960)
	output := make([]byte, 4000)
	result := opusEncode(12345, unsafe.Pointer(&pcm[0]), 960, unsafe.Pointer(&output[0]), 4000)
	if result != 960 {
		t.Errorf("opusEncode() = %d, want 960", result)
	}
}

func TestProcessData(t *testing.T) {
	a := int32(10)
	b := byte(20)
	result := processData(&a, &b, 30)
	if result != 60 {
		t.Errorf("processData(&10, &20, 30) = %d, want 60", result)
	}
}

func TestNoReturnFn(t *testing.T) {
	// Just verify it doesn't panic
	noReturnFn(1, 2, 3)
}

func TestManyArgs(t *testing.T) {
	result := manyArgs(1, 2, 3, 4, 5, 6, 7, 8, 9)
	expected := int64(1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9)
	if result != expected {
		t.Errorf("manyArgs(1..9) = %d, want %d", result, expected)
	}
}

func TestBoolAnd(t *testing.T) {
	tests := []struct {
		a, b, want bool
	}{
		{true, true, true},
		{true, false, false},
		{false, true, false},
		{false, false, false},
	}

	for _, tt := range tests {
		got := boolAnd(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("boolAnd(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

// Benchmark tests

func BenchmarkBenchAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchAdd(10, 20)
	}
}

func BenchmarkBenchAddDirect(b *testing.B) {
	sym, _ := purego.Dlsym(testLib, "bench_add")
	for i := 0; i < b.N; i++ {
		_ = purego.Syscall2(sym, 10, 20)
	}
}

func BenchmarkBenchPtr(b *testing.B) {
	data := make([]byte, 100)
	ptr := unsafe.Pointer(&data[0])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchPtr(ptr, 50)
	}
}

func BenchmarkOpusEncode(b *testing.B) {
	pcm := make([]int16, 960)
	output := make([]byte, 4000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opusEncode(12345, unsafe.Pointer(&pcm[0]), 960, unsafe.Pointer(&output[0]), 4000)
	}
}

func BenchmarkIsValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = isValid(int32(i & 1))
	}
}

func BenchmarkManyArgs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = manyArgs(1, 2, 3, 4, 5, 6, 7, 8, 9)
	}
}
