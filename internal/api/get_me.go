// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) GetMe(
	ctx context.Context,
	req GetMeRequestObject,
) (GetMeResponseObject, error) {
	// The gate put the actor here, having read their token and found the
	// session behind it still open. Absence means this handler was reached
	// without the gate in front of it, which is a fault of ours.
	actor, present := actorFrom(ctx)
	if !present {
		return GetMe500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	info, getErr := usecase.NewGetMe(s.pool, actor.ID).Exec(ctx)
	if getErr != nil {
		return getMeError(getErr), nil
	}

	user, err := newUser(info)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return GetMe500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	return GetMe200JSONResponse(user), nil
}

// getMeError is a type switch rather than a chain of checks so that
// gochecksumtype verifies it: a failure added to [domain.GetMeError]
// breaks the build here instead of quietly becoming a 500.
func getMeError(err domain.GetMeError) GetMeResponseObject {
	switch err.(type) {
	case domain.UserNotFoundError:
		return GetMe401ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusUnauthorized),
		)
	case domain.UnexpectedError:
		return GetMe500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return GetMe500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}
