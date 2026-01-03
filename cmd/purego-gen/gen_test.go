// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Ebitengine Authors

package main

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenClosure tests closure generation for various function signatures
func TestGenClosure(t *testing.T) {
	tests := []struct {
		name     string
		fi       funcInfo
		contains []string // strings that should be in output
	}{
		{
			name: "simple_int",
			fi: funcInfo{
				VarName: "foo",
				LibVar:  "lib",
				SymName: "foo",
				Args:    []argInfo{{Name: "a0", Type: "int32"}},
				RetType: "int32",
				HasRet:  true,
			},
			contains: []string{"func(a0 int32) int32", "Syscall1", "uintptr(a0)", "return int32"},
		},
		{
			name: "no_return",
			fi: funcInfo{
				VarName: "bar",
				LibVar:  "lib",
				SymName: "bar",
				Args:    []argInfo{{Name: "a0", Type: "int64"}, {Name: "a1", Type: "int64"}},
				HasRet:  false,
			},
			contains: []string{"func(a0 int64, a1 int64)", "Syscall2"},
		},
		{
			name: "unsafe_pointer_keepalive",
			fi: funcInfo{
				VarName: "encode",
				LibVar:  "lib",
				SymName: "encode",
				Args: []argInfo{
					{Name: "a0", Type: "uint64"},
					{Name: "a1", Type: "unsafe.Pointer"},
					{Name: "a2", Type: "int32"},
				},
				RetType: "int32",
				HasRet:  true,
			},
			contains: []string{"runtime.KeepAlive(a1)", "_r :=", "return _r"},
		},
		{
			name: "typed_pointer_keepalive",
			fi: funcInfo{
				VarName: "process",
				LibVar:  "lib",
				SymName: "process",
				Args: []argInfo{
					{Name: "a0", Type: "*int32"},
					{Name: "a1", Type: "*byte"},
				},
				RetType: "int32",
				HasRet:  true,
			},
			contains: []string{"runtime.KeepAlive(a0)", "runtime.KeepAlive(a1)", "uintptr(unsafe.Pointer(a0))"},
		},
		{
			name: "nine_args_syscall9",
			fi: funcInfo{
				VarName: "big",
				LibVar:  "lib",
				SymName: "big",
				Args: []argInfo{
					{Name: "a0", Type: "int64"},
					{Name: "a1", Type: "int64"},
					{Name: "a2", Type: "int64"},
					{Name: "a3", Type: "int64"},
					{Name: "a4", Type: "int64"},
					{Name: "a5", Type: "int64"},
					{Name: "a6", Type: "int64"},
					{Name: "a7", Type: "int64"},
					{Name: "a8", Type: "int64"},
				},
				RetType: "int64",
				HasRet:  true,
			},
			contains: []string{"Syscall9"}, // Now uses Syscall9 instead of SyscallN
		},
		{
			name: "uintptr_passthrough",
			fi: funcInfo{
				VarName: "ptr",
				LibVar:  "lib",
				SymName: "ptr",
				Args:    []argInfo{{Name: "a0", Type: "uintptr"}},
				RetType: "uintptr",
				HasRet:  true,
			},
			contains: []string{"a0)"}, // should NOT wrap in uintptr()
		},
		{
			name: "zero_args",
			fi: funcInfo{
				VarName: "nop",
				LibVar:  "lib",
				SymName: "nop",
				Args:    []argInfo{},
				HasRet:  false,
			},
			contains: []string{"func()", "Syscall0"},
		},
		{
			name: "zero_args_with_return",
			fi: funcInfo{
				VarName: "getval",
				LibVar:  "lib",
				SymName: "getval",
				Args:    []argInfo{},
				RetType: "int64",
				HasRet:  true,
			},
			contains: []string{"func() int64", "Syscall0", "return int64"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := genClosure(tt.fi)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("genClosure() missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

// TestConv tests argument conversion
func TestConv(t *testing.T) {
	tests := []struct {
		name, typ, want string
	}{
		{"int", "int32", "uintptr(a)"},
		{"uintptr", "uintptr", "a"},
		{"unsafe_pointer", "unsafe.Pointer", "uintptr(a)"},
		{"typed_pointer", "*int32", "uintptr(unsafe.Pointer(a))"},
		{"double_pointer", "**byte", "uintptr(unsafe.Pointer(a))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := conv("a", tt.typ)
			if got != tt.want {
				t.Errorf("conv(%q, %q) = %q, want %q", "a", tt.typ, got, tt.want)
			}
		})
	}
}

// TestNeedsKeepAlive tests keep-alive detection
func TestNeedsKeepAlive(t *testing.T) {
	tests := []struct {
		typ  string
		want bool
	}{
		{"int32", false},
		{"int64", false},
		{"uintptr", false},
		{"unsafe.Pointer", true},
		{"*int32", true},
		{"*byte", true},
		{"**int", true},
	}

	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			got := needsKeepAlive(tt.typ)
			if got != tt.want {
				t.Errorf("needsKeepAlive(%q) = %v, want %v", tt.typ, got, tt.want)
			}
		})
	}
}

// TestGeneratedCodeFormatsCorrectly verifies generated code is valid Go syntax
func TestGeneratedCodeFormatsCorrectly(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name: "basic",
			input: `package test
import "unsafe"
var lib uintptr
//purego:sym lib simple
var simple func(int32) int32
`,
		},
		{
			name: "with_pointer",
			input: `package test
import "unsafe"
var lib uintptr
//purego:sym lib encode
var encode func(uint64, unsafe.Pointer, int32) int32
`,
		},
		{
			name: "typed_pointer",
			input: `package test
import "unsafe"
var lib uintptr
//purego:sym lib process
var process func(*int32, *byte) int32
`,
		},
		{
			name: "no_return",
			input: `package test
import "unsafe"
var lib uintptr
//purego:sym lib callback
var callback func(uint64, unsafe.Pointer)
`,
		},
		{
			name: "many_functions",
			input: `package test
import "unsafe"
var lib uintptr
//purego:sym lib fn1
var fn1 func(int32) int32
//purego:sym lib fn2
var fn2 func(int64, int64) int64
//purego:sym lib fn3
var fn3 func(unsafe.Pointer) unsafe.Pointer
`,
		},
		{
			name: "many_args",
			input: `package test
import "unsafe"
var lib uintptr
//purego:sym lib bigfn
var bigfn func(int64, int64, int64, int64, int64, int64, int64, int64, int64) int64
`,
		},
		{
			name: "complex_mixed_pointers",
			input: `package test
import "unsafe"
var lib uintptr
//purego:sym lib complex
var complex func(*int32, unsafe.Pointer, **byte, int64, *uint64) int32
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()

			// Write input file
			inputPath := filepath.Join(dir, "input.go")
			if err := os.WriteFile(inputPath, []byte(tc.input), 0644); err != nil {
				t.Fatal(err)
			}

			// Run generator
			oldOutput := *output
			*output = filepath.Join(dir, "purego_gen.go")
			defer func() { *output = oldOutput }()

			// Change to temp dir and run
			oldWd, _ := os.Getwd()
			os.Chdir(dir)
			defer os.Chdir(oldWd)

			main() // This will parse and generate

			// Read generated file
			generated, err := os.ReadFile(*output)
			if err != nil {
				t.Fatalf("Generator didn't create output: %v", err)
			}

			// Check it's valid Go syntax
			if _, err := format.Source(generated); err != nil {
				t.Fatalf("Generated invalid Go code: %v\n%s", err, string(generated))
			}
		})
	}
}

// TestEdgeCases tests that the generator handles edge cases gracefully
func TestEdgeCases(t *testing.T) {
	// These should NOT panic - generator should handle gracefully
	edgeCases := []funcInfo{
		// Empty name edge case
		{VarName: "", LibVar: "lib", SymName: "sym", Args: nil, HasRet: false},
		// Weird type names
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: "MyCustomType"}}, HasRet: true, RetType: "MyResult"},
		// Slice types (not supported but shouldn't crash)
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: "[]int32"}}, HasRet: false},
		// Map type (not supported but shouldn't crash)
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: "map[string]int"}}, HasRet: false},
		// Channel type (not supported but shouldn't crash)
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: "chan int"}}, HasRet: false},
		// Interface type
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: "any"}}, HasRet: true, RetType: "error"},
		// Empty string type
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: ""}}, HasRet: false},
		// Very long type name
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: "VeryLongTypeNameThatShouldStillWork123456789"}}, HasRet: false},
		// Complex nested pointer
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: "***int32"}}, HasRet: true, RetType: "***byte"},
		// Pointer to unsafe.Pointer
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{{Name: "a0", Type: "*unsafe.Pointer"}}, HasRet: false},
		// Many identical args
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: func() []argInfo {
				var a []argInfo
				for i := 0; i < 20; i++ {
					a = append(a, argInfo{Name: "a" + string(rune('0'+i%10)), Type: "int64"})
				}
				return a
			}(), HasRet: true, RetType: "int64"},
		// All pointer args
		{VarName: "fn", LibVar: "lib", SymName: "sym",
			Args: []argInfo{
				{Name: "a", Type: "unsafe.Pointer"},
				{Name: "b", Type: "*int32"},
				{Name: "c", Type: "*byte"},
				{Name: "d", Type: "**int"},
				{Name: "e", Type: "*uint64"},
			}, HasRet: true, RetType: "*int32"},
	}

	for i, tc := range edgeCases {
		t.Run(strings.ReplaceAll(tc.VarName+"_"+tc.SymName+"_"+string(rune('0'+i)), " ", "_"), func(t *testing.T) {
			// Should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("genClosure panicked: %v", r)
				}
			}()
			_ = genClosure(tc)
		})
	}
}

// TestExprToStringEdgeCases tests AST expression parsing
func TestExprToStringEdgeCases(t *testing.T) {
	// Test that the function handles nil gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("exprToString panicked on nil: %v", r)
		}
	}()
	// nil input should return "any"
	result := exprToString(nil)
	if result != "any" {
		t.Errorf("exprToString(nil) = %q, want %q", result, "any")
	}
}

// FuzzFuncSignature fuzzes different function signatures
func FuzzFuncSignature(f *testing.F) {
	// Seed with basic cases
	f.Add(uint8(1), uint8(1), true, true) // 1 arg, 1 pointer, has return, has keepalive
	f.Add(uint8(0), uint8(0), true, false)
	f.Add(uint8(3), uint8(0), false, false)
	f.Add(uint8(5), uint8(2), true, true)
	f.Add(uint8(8), uint8(0), true, false)
	f.Add(uint8(9), uint8(5), true, true) // triggers SyscallN

	types := []string{"int32", "int64", "uint32", "uint64", "uintptr"}
	ptrTypes := []string{"unsafe.Pointer", "*int32", "*byte"}

	f.Fuzz(func(t *testing.T, numArgs uint8, numPtrs uint8, hasReturn bool, needsKeep bool) {
		// Limit to reasonable sizes
		numArgs = numArgs % 12
		if numPtrs > numArgs {
			numPtrs = numArgs
		}

		var args []argInfo
		for i := uint8(0); i < numArgs; i++ {
			var typ string
			if i < numPtrs {
				typ = ptrTypes[int(i)%len(ptrTypes)]
			} else {
				typ = types[int(i)%len(types)]
			}
			args = append(args, argInfo{
				Name: "a" + string('0'+i%10),
				Type: typ,
			})
		}

		fi := funcInfo{
			VarName: "fn",
			LibVar:  "lib",
			SymName: "fn",
			Args:    args,
			HasRet:  hasReturn,
		}
		if hasReturn {
			fi.RetType = "int64"
		}

		// Should not panic
		result := genClosure(fi)

		// Basic sanity checks
		if !strings.Contains(result, "func(") {
			t.Error("missing func(")
		}
		if hasReturn && !strings.Contains(result, "return") {
			t.Error("missing return for function with return value")
		}

		// Check KeepAlive for pointer args
		for _, arg := range args {
			if needsKeepAlive(arg.Type) {
				if !strings.Contains(result, "runtime.KeepAlive("+arg.Name+")") {
					t.Errorf("missing KeepAlive for %s (type %s)", arg.Name, arg.Type)
				}
			}
		}

		// Check proper syscall variant (now we have Syscall0-15)
		if len(args) > 15 && !strings.Contains(result, "SyscallN") {
			t.Error("should use SyscallN for >15 args")
		}
		if len(args) <= 15 && strings.Contains(result, "SyscallN") {
			t.Error("should not use SyscallN for <=15 args")
		}
	})
}

// FuzzRandomTypes tests with random type strings
func FuzzRandomTypes(f *testing.F) {
	f.Add("int32", "int64", true)
	f.Add("*byte", "unsafe.Pointer", false)
	f.Add("uintptr", "uint64", true)
	f.Add("", "any", true)
	f.Add("***int", "*unsafe.Pointer", false)

	f.Fuzz(func(t *testing.T, argType, retType string, hasRet bool) {
		fi := funcInfo{
			VarName: "fn",
			LibVar:  "lib",
			SymName: "fn",
			Args:    []argInfo{{Name: "a0", Type: argType}},
			HasRet:  hasRet,
			RetType: retType,
		}

		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("genClosure panicked with argType=%q retType=%q: %v", argType, retType, r)
			}
		}()

		result := genClosure(fi)
		if !strings.Contains(result, "func(") {
			t.Error("missing func(")
		}
	})
}
