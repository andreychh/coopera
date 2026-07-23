// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"encoding/json"
	"net/http"
)

// RequestError is a StrictHTTPServerOptions.RequestErrorHandlerFunc. It
// runs for requests the strict handler rejects before Server sees them,
// e.g. a body that isn't valid JSON or that doesn't match the schema's
// types (a number where a string is expected).
func RequestError(w http.ResponseWriter, _ *http.Request, _ error) {
	writeProblem(w, http.StatusBadRequest, "request body is not valid JSON")
}

// ResponseError is a StrictHTTPServerOptions.ResponseErrorHandlerFunc. It
// runs when a Server response can't be encoded.
func ResponseError(w http.ResponseWriter, _ *http.Request, _ error) {
	writeProblem(w, http.StatusInternalServerError, "")
}

// NewProblem builds a Problem for status, with no further detail.
func NewProblem(status int) Problem {
	return Problem{
		Title:  http.StatusText(status),
		Status: status,
	}
}

// NewDetailedProblem builds a Problem for status, with detail explaining
// what specifically went wrong.
func NewDetailedProblem(status int, detail string) Problem {
	problem := NewProblem(status)
	problem.Detail = new(detail)
	return problem
}

func writeProblem(w http.ResponseWriter, status int, detail string) {
	problem := NewProblem(status)
	if detail != "" {
		problem = NewDetailedProblem(status, detail)
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	//nolint:errchkjson // status and headers are already sent; nothing left to do if this fails
	_ = json.NewEncoder(w).Encode(problem)
}
