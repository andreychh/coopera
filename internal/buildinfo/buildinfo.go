// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package buildinfo provides access to binary metadata set during the build
// process.
//
// Most of the information is typically injected via linker flags (ldflags).
package buildinfo

import (
	"runtime"
)

// The following variables are intended to be populated at build time.
//
//nolint:gochecknoglobals // These variables are set at build time via ldflags
var (
	version   = "dev"
	commit    = "unknown"
	branch    = "unknown"
	buildTime = "unknown"
	treeState = "unknown"
	builtBy   = "manual"
)

// BuildInfo holds metadata about the compiled binary.
type BuildInfo struct {
	// Version is the semantic version of the build (e.g., "v1.2.3").
	Version string
	// Commit is the full git SHA of the commit used for the build.
	Commit string
	// Branch is the name of the git branch.
	Branch string
	// BuildTime is the RFC3339 formatted timestamp of the build.
	BuildTime string
	// TreeState indicates if the git working tree was clean or dirty.
	TreeState string
	// BuiltBy identifies the tool or environment that triggered the build.
	BuiltBy string
	// GoVersion is the version of the Go compiler used.
	GoVersion string
	// Platform is the target OS and Architecture (e.g., "linux/amd64").
	Platform string
}

// Read gathers the build-time metadata and current runtime information into a
// single BuildInfo structure.
func Read() *BuildInfo {
	return &BuildInfo{
		Version:   version,
		Commit:    commit,
		Branch:    branch,
		BuildTime: buildTime,
		BuiltBy:   builtBy,
		TreeState: treeState,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}
}
