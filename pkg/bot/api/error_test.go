// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/andreychh/coopera/pkg/bot/api"
	"github.com/andreychh/coopera/pkg/ptr"
)

func TestError_Is(t *testing.T) {
	apiErr := api.NewError(api.Envelope{
		ErrorCode:   ptr.Ptr[int32](404),
		Description: ptr.Ptr("Not Found"),
		Parameters:  nil,
		Result:      nil,
		Ok:          false,
	})
	wrappedErr := fmt.Errorf("outer wrap: %w", apiErr)
	if !errors.Is(wrappedErr, apiErr) {
		t.Errorf("expected errors.Is to find the API error instance in the chain")
	}
}

func TestError_As(t *testing.T) {
	const expectedCode int32 = 404
	apiErr := api.NewError(api.Envelope{
		ErrorCode:   ptr.Ptr(expectedCode),
		Description: ptr.Ptr("Not Found"),
		Parameters:  nil,
		Result:      nil,
		Ok:          false,
	})
	wrappedErr := fmt.Errorf("outer wrap: %w", apiErr)

	var target *api.Error
	if !errors.As(wrappedErr, &target) {
		t.Fatalf("expected errors.As to successfully extract *api.Error from the chain")
	}
	if *target.Envelope().ErrorCode != expectedCode {
		t.Errorf("expected error code %d, got %d", expectedCode, *target.Envelope().ErrorCode)
	}
}
