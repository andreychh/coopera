// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"context"
	"net/http"

	"github.com/andreychh/coopera/internal/domain"
	"github.com/andreychh/coopera/internal/usecase"
)

func (s Server) CreateCode(
	ctx context.Context,
	req CreateCodeRequestObject,
) (CreateCodeResponseObject, error) {
	email, err := domain.ParseEmail(string(req.Body.Email))
	if err != nil {
		//nolint:nilerr // outcome is encoded in the response, not the error return
		return CreateCode400ApplicationProblemPlusJSONResponse(
			NewDetailedProblem(http.StatusBadRequest, "Invalid email"),
		), nil
	}

	// Kept in its own variable: assigning into the err above would widen
	// it back to error and the sum would stop being checked.
	delivery, createErr := usecase.NewCreateCode(s.pool, s.post, email).Exec(ctx)
	if createErr != nil {
		return createCodeError(createErr), nil
	}

	return CreateCode202JSONResponse(newCodeDelivery(delivery)), nil
}

// createCodeError is a type switch rather than a chain of checks so that
// gochecksumtype verifies it: a failure added to [domain.CreateCodeError]
// breaks the build here instead of quietly becoming a 500.
func createCodeError(err domain.CreateCodeError) CreateCodeResponseObject {
	switch e := err.(type) {
	case domain.TooManyCodesError:
		return CreateCode429ApplicationProblemPlusJSONResponse{
			Body: NewDetailedProblem(
				http.StatusTooManyRequests,
				"Too many codes requested for this address",
			),
			Headers: CreateCode429ResponseHeaders{
				RetryAfter: new(seconds(e.RetryAfter)),
			},
		}
	case domain.UnexpectedError:
		return CreateCode500ApplicationProblemPlusJSONResponse(
			NewProblem(http.StatusInternalServerError),
		)
	}
	return CreateCode500ApplicationProblemPlusJSONResponse(
		NewProblem(http.StatusInternalServerError),
	)
}
