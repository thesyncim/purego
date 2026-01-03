// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build (darwin || freebsd || linux || netbsd) && (amd64 || arm64)

package purego

import (
	"unsafe"
)

// Syscall0 through Syscall8 are specialized versions of SyscallN that avoid
// variadic overhead and array copying. Zero allocations when the pool is warm.

//go:nosplit
func Syscall0(fn uintptr) uintptr {
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
