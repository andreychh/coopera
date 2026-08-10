// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"
)

// The operations below are described in openapi.yml but not built yet.
// They answer 500 so the generated interface is satisfied and the tree
// compiles; each disappears from here the moment its own file appears.

func (s Server) CreateSession(
	_ context.Context,
	_ CreateSessionRequestObject,
) (CreateSessionResponseObject, error) {
	return CreateSession500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	), nil
}

func (s Server) RefreshSession(
	_ context.Context,
	_ RefreshSessionRequestObject,
) (RefreshSessionResponseObject, error) {
	return RefreshSession500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	), nil
}

func (s Server) DeleteSession(
	_ context.Context,
	_ DeleteSessionRequestObject,
) (DeleteSessionResponseObject, error) {
	return DeleteSession500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	), nil
}
