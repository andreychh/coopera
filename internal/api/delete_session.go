// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) DeleteSession(
	ctx context.Context,
	_ DeleteSessionRequestObject,
) (DeleteSessionResponseObject, error) {
	// The gate put the actor here, having read their token and found the
	// session behind it still open. Absence means this handler was reached
	// without the gate in front of it, which is a fault of ours.
	actor, present := actorFrom(ctx)
	if !present {
		return DeleteSession500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	// Kept in its own variable: assigning into an err would widen it back
	// to error and the sum would stop being checked.
	endErr := usecase.NewEndSession(s.pool, actor).Exec(ctx)
	if endErr != nil {
		return deleteSessionError(endErr), nil
	}

	return DeleteSession204Response{}, nil
}

// deleteSessionError is a type switch rather than a chain of checks so
// that gochecksumtype verifies it: a failure added to
// [domain.EndSessionError] breaks the build here instead of quietly
// becoming a 500.
func deleteSessionError(err domain.EndSessionError) DeleteSessionResponseObject {
	//nolint:gocritic // the switch is what gochecksumtype checks; an if is not
	switch err.(type) {
	case domain.UnexpectedError:
		return DeleteSession500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return DeleteSession500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}
