// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) CreateSession(
	ctx context.Context,
	req CreateSessionRequestObject,
) (CreateSessionResponseObject, error) {
	email, err := domain.ParseEmail(string(req.Body.Email))
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateSession400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid email"),
		), nil
	}

	code, err := domain.ParseSignInCode(req.Body.Code)
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateSession400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid code"),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	pass, createErr := usecase.NewCreateSession(s.pool, s.seal, email, code).Exec(ctx)
	if createErr != nil {
		return createSessionError(createErr), nil
	}

	// Every session answers from the same address: it is reached through
	// the pass that holds it, so there is nothing in a URL to tell one
	// session from another.
	return CreateSession201JSONResponse{
		Body: newPass(pass),
		Headers: CreateSession201ResponseHeaders{
			Location: new("/v1/auth/sessions/current"),
		},
	}, nil
}

// createSessionError is a type switch rather than a chain of checks so
// that gochecksumtype verifies it: a failure added to
// [domain.CreateSessionError] breaks the build here instead of quietly
// becoming a 500.
func createSessionError(err domain.CreateSessionError) CreateSessionResponseObject {
	switch e := err.(type) {
	case domain.SignInCodeMismatchError:
		return CreateSession401ApplicationProblemPlusJSONResponse{
			AttemptsLeft: int(e.AttemptsLeft),
			Detail:       new("The code does not match"),
			Status:       http.StatusUnauthorized,
			Title:        http.StatusText(http.StatusUnauthorized),
		}
	case domain.SignInCodeNotUsableError:
		return CreateSession410ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusGone, "Ask for a new code"),
		)
	case domain.UnexpectedError:
		return CreateSession500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return CreateSession500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}
