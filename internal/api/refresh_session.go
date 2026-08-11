// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) RefreshSession(
	ctx context.Context,
	req RefreshSessionRequestObject,
) (RefreshSessionResponseObject, error) {
	// The key is taken as given, without a look at its shape. There is
	// nothing a shape could protect here: an unknown key costs nothing to
	// refuse, and there is no count for a wrong one to burn.
	token := domain.RefreshToken(req.Body.RefreshToken)

	// Kept in its own variable: assigning into an err would widen it back
	// to error and the sum would stop being checked.
	pass, refreshErr := usecase.NewRefreshSession(s.pool, s.seal, token).Exec(ctx)
	if refreshErr != nil {
		return refreshSessionError(refreshErr), nil
	}

	return RefreshSession200JSONResponse(newPass(pass)), nil
}

// refreshSessionError is a type switch rather than a chain of checks so
// that gochecksumtype verifies it: a failure added to
// [domain.RefreshSessionError] breaks the build here instead of quietly
// becoming a 500.
func refreshSessionError(err domain.RefreshSessionError) RefreshSessionResponseObject {
	switch err.(type) {
	case domain.RefreshTokenNotUsableError:
		return RefreshSession401ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusUnauthorized, "Sign in again"),
		)
	case domain.UnexpectedError:
		return RefreshSession500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return RefreshSession500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}
