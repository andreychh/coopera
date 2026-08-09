// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

// ForeignKeyViolation is the Postgres SQLSTATE code for a violated
// foreign key (https://www.postgresql.org/docs/current/errcodes-appendix.html),
// and MembersUserIDForeignKey is the name Postgres generates for
// members' unnamed "user_id UUID NOT NULL REFERENCES users (id)"
// constraint, following its default <table>_<column>_fkey convention.
//
// UniqueViolation is the SQLSTATE code for a value that is already
// taken. Unlike the pair above it needs no constraint name: the only
// statement that raises it here updates a single unique column, so the
// code alone says which value clashed. Should another unique column ever
// be written by the same statement, this shortcut stops being true.
//
// This file is meant to look out of place. It is where the domain speaks
// Postgres, and the constraint name is the brittle part: renaming it in
// a migration would break the match without a word from the compiler.
// That pair survives only because nothing but the database can tell that
// an actor does not exist, and it goes once authentication settles that
// question ahead of every request. Every other such translation has
// already been replaced by a query whose result carries the answer.
const (
	ForeignKeyViolation     = "23503"
	MembersUserIDForeignKey = "members_user_id_fkey"
	UniqueViolation         = "23505"
)
