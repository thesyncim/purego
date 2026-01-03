// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 The Ebitengine Authors

//go:build !(darwin || freebsd || linux || netbsd) || !(amd64 || arm64)

package purego

// tryRegisterTyped returns false on unsupported platforms.
// The optimized typed path requires amd64 or arm64 on unix-like systems.
func tryRegisterTyped(fptr any, cfn uintptr) bool {
	return false
}
