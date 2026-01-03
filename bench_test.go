// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Ebitengine Authors

package purego_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
)

var (
	benchLib    uintptr
	benchNoop   func() int64
	bench1Int   func(int64) int64
	bench3Int   func(int64, int64, int64) int64
	bench6Int   func(int64, int64, int64, int64, int64, int64) int64
	bench8Int   func(int64, int64, int64, int64, int64, int64, int64, int64) int64
	bench1Float func(float64) float64
	bench3Float func(float64, float64, float64) float64
	benchMixed  func(int64, float64, int64, float64) float64

	// Direct syscall function pointers (for baseline comparison)
	benchNoopSym uintptr
	bench1IntSym uintptr
	bench3IntSym uintptr
	bench6IntSym uintptr
	bench8IntSym uintptr
)

func setupBenchLib(b *testing.B) {
	if benchLib != 0 {
		return
	}

	libFileName := filepath.Join(b.TempDir(), "benchtest.so")
	if err := buildSharedLib("CC", libFileName, filepath.Join("testdata", "abitest", "abi_test.c")); err != nil {
		b.Fatal(err)
	}

	var err error
	benchLib, err = load.OpenLibrary(libFileName)
	if err != nil {
		b.Fatalf("Failed to open library: %v", err)
	}

	// RegisterLibFunc now uses optimized typed path for common signatures
	purego.RegisterLibFunc(&benchNoop, benchLib, "bench_noop")
	purego.RegisterLibFunc(&bench1Int, benchLib, "bench_1int")
	purego.RegisterLibFunc(&bench3Int, benchLib, "bench_3int")
	purego.RegisterLibFunc(&bench6Int, benchLib, "bench_6int")
	purego.RegisterLibFunc(&bench8Int, benchLib, "bench_8int")
	purego.RegisterLibFunc(&bench1Float, benchLib, "bench_1float")
	purego.RegisterLibFunc(&bench3Float, benchLib, "bench_3float")
	purego.RegisterLibFunc(&benchMixed, benchLib, "bench_mixed")

	// Direct syscall symbols for raw performance comparison
	benchNoopSym, _ = purego.Dlsym(benchLib, "bench_noop")
	bench1IntSym, _ = purego.Dlsym(benchLib, "bench_1int")
	bench3IntSym, _ = purego.Dlsym(benchLib, "bench_3int")
	bench6IntSym, _ = purego.Dlsym(benchLib, "bench_6int")
	bench8IntSym, _ = purego.Dlsym(benchLib, "bench_8int")
}

// RegisterFunc benchmarks - these now use optimized typed paths internally
// for common signatures (zero allocations)

func BenchmarkCallNoop(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchNoop()
	}
}

func BenchmarkCall1Int(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bench1Int(42)
	}
}

func BenchmarkCall3Int(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bench3Int(1, 2, 3)
	}
}

func BenchmarkCall6Int(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bench6Int(1, 2, 3, 4, 5, 6)
	}
}

func BenchmarkCall8Int(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bench8Int(1, 2, 3, 4, 5, 6, 7, 8)
	}
}

func BenchmarkCall1Float(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bench1Float(3.14)
	}
}

func BenchmarkCall3Float(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bench3Float(1.0, 2.0, 3.0)
	}
}

func BenchmarkCallMixed(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchMixed(1, 2.0, 3, 4.0)
	}
}

// Direct Syscall benchmarks - raw performance baseline (zero allocations)

func BenchmarkSyscall0(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = purego.Syscall0(benchNoopSym)
	}
}

func BenchmarkSyscall1(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = purego.Syscall1(bench1IntSym, 42)
	}
}

func BenchmarkSyscall3(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = purego.Syscall3(bench3IntSym, 1, 2, 3)
	}
}

func BenchmarkSyscall6(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = purego.Syscall6(bench6IntSym, 1, 2, 3, 4, 5, 6)
	}
}

func BenchmarkSyscall8(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = purego.Syscall8(bench8IntSym, 1, 2, 3, 4, 5, 6, 7, 8)
	}
}

func BenchmarkSyscallN_8args(b *testing.B) {
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		b.Skip("benchmark requires amd64 or arm64")
	}
	setupBenchLib(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = purego.SyscallN(bench8IntSym, 1, 2, 3, 4, 5, 6, 7, 8)
	}
}
