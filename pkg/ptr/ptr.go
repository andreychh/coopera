// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ptr

// Zero returns the zero value of type T. Facilitates initializing generic types
// within functions.
//
//nolint:ireturn // Returns a generic type parameter.
func Zero[T any]() T {
	var zero T
	return zero
}

// ValueOrZero returns the value pointed to by ptr if it's not nil. If ptr is
// nil, it returns the zero value for type T.
//
//nolint:ireturn // Returns a generic type parameter.
func ValueOrZero[T any](ptr *T) T {
	if ptr == nil {
		return Zero[T]()
	}
	return *ptr
}

// ValueOrDefault returns the value pointed to by ptr if it's not nil. If ptr is
// nil, it returns the provided default value.
//
//nolint:ireturn // Returns a generic type parameter.
func ValueOrDefault[T any](ptr *T, def T) T {
	if ptr == nil {
		return def
	}
	return *ptr
}
