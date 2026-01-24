// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package utils

func Ptr[T any](value T) *T {
	return &value
}

func Zero[T any]() T {
	var zero T
	return zero
}

func ValueOrZero[T any](ptr *T) T {
	if ptr == nil {
		return Zero[T]()
	}
	return *ptr
}

func ValueOrDefault[T any](ptr *T, def T) T {
	if ptr == nil {
		return def
	}
	return *ptr
}
