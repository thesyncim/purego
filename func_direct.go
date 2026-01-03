// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build (darwin || freebsd || linux || netbsd) && (amd64 || arm64)

package purego

import (
	"unsafe"
)

// Syscall0 through Syscall15 are specialized versions of SyscallN that avoid
// variadic overhead and array copying. Zero allocations when the pool is warm.

//go:nosplit
func Syscall0(fn uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall1(fn, a1 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall2(fn, a1, a2 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall3(fn, a1, a2, a3 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall4(fn, a1, a2, a3, a4 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall5(fn, a1, a2, a3, a4, a5 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall7(fn, a1, a2, a3, a4, a5, a6, a7 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall8(fn, a1, a2, a3, a4, a5, a6, a7, a8 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	syscall.a8 = a8
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall9(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	syscall.a8 = a8
	syscall.a9 = a9
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall10(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	syscall.a8 = a8
	syscall.a9 = a9
	syscall.a10 = a10
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall11(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	syscall.a8 = a8
	syscall.a9 = a9
	syscall.a10 = a10
	syscall.a11 = a11
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall12(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	syscall.a8 = a8
	syscall.a9 = a9
	syscall.a10 = a10
	syscall.a11 = a11
	syscall.a12 = a12
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall13(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	syscall.a8 = a8
	syscall.a9 = a9
	syscall.a10 = a10
	syscall.a11 = a11
	syscall.a12 = a12
	syscall.a13 = a13
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall14(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	syscall.a8 = a8
	syscall.a9 = a9
	syscall.a10 = a10
	syscall.a11 = a11
	syscall.a12 = a12
	syscall.a13 = a13
	syscall.a14 = a14
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}

//go:nosplit
func Syscall15(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) uintptr {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	syscall := thePool.Get().(*syscall15Args)
	syscall.fn = fn
	syscall.a1 = a1
	syscall.a2 = a2
	syscall.a3 = a3
	syscall.a4 = a4
	syscall.a5 = a5
	syscall.a6 = a6
	syscall.a7 = a7
	syscall.a8 = a8
	syscall.a9 = a9
	syscall.a10 = a10
	syscall.a11 = a11
	syscall.a12 = a12
	syscall.a13 = a13
	syscall.a14 = a14
	syscall.a15 = a15
	runtime_cgocall(syscall15XABI0, unsafe.Pointer(syscall))
	r := syscall.a1
	*syscall = syscall15Args{}
	thePool.Put(syscall)
	return r
}
