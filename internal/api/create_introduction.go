// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) CreateIntroduction(
	ctx context.Context,
	req CreateIntroductionRequestObject,
) (CreateIntroductionResponseObject, error) {
	username, err := domain.ParseUsername(req.Body.Username)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateIntroduction400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid username"),
		), nil
	}
	userID, err := domain.ParseID(req.Params.XUserId)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateIntroduction400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid X-User-Id"),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	info, createErr := usecase.NewCreateIntroduction(s.pool, userID, username).Exec(ctx)
	if createErr != nil {
		return createIntroductionError(createErr), nil
	}

	user, err := newUser(info)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateIntroduction500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		), nil
	}

	return CreateIntroduction201JSONResponse(user), nil
}

// createIntroductionError is a type switch rather than a chain of checks
// so that gochecksumtype verifies it: a failure added to
// [domain.CreateIntroductionError] breaks the build here instead of
// quietly becoming a 500.
//
// Both conflicts answer 409 and differ only in wording, which is all the
// difference a client can see until problems carry a type.
func createIntroductionError(err domain.CreateIntroductionError) CreateIntroductionResponseObject {
	switch err.(type) {
	case domain.UserNotFoundError:
		return CreateIntroduction401ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusUnauthorized),
		)
	case domain.AlreadyIntroducedError:
		return CreateIntroduction409ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusConflict, "Already introduced"),
		)
	case domain.UsernameTakenError:
		return CreateIntroduction409ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusConflict, "Username is taken"),
		)
	case domain.UnexpectedError:
		return CreateIntroduction500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return CreateIntroduction500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}
