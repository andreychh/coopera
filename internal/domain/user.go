// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package domain

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/andreychh/coopera/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// foreignKeyViolation is the Postgres SQLSTATE code for a foreign key
// constraint violation (https://www.postgresql.org/docs/current/errcodes-appendix.html).
const foreignKeyViolation = "23503"

// membersUserIDForeignKey is the name Postgres generates for the
// unnamed "user_id UUID NOT NULL REFERENCES users (id)" constraint on
// members, following its default <table>_<column>_fkey convention.
const membersUserIDForeignKey = "members_user_id_fkey"

type ID uuid.UUID

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}

type TeamName string

func ParseTeamName(s string) (TeamName, error) {
	if strings.TrimSpace(s) != s {
		return "", errors.New("must not have leading or trailing whitespace")
	}
	count := utf8.RuneCountInString(s)
	if count < 1 || count > 100 {
		return "", errors.New("must be between 1 and 100 characters")
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return "", errors.New("must not contain control characters")
		}
	}
	return TeamName(s), nil
}

func (n TeamName) String() string {
	return string(n)
}

type DateTime time.Time

func ParseDateTime(s string) (DateTime, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return DateTime{}, fmt.Errorf("invalid format: %w", err)
	}
	return DateTime(t), nil
}

type InviteLinkExpiry DateTime

func ParseInviteLinkExpiry(s string) (InviteLinkExpiry, error) {
	dt, err := ParseDateTime(s)
	if err != nil {
		return InviteLinkExpiry{}, err
	}
	if !time.Time(dt).After(time.Now()) {
		return InviteLinkExpiry{}, errors.New("must be in the future")
	}
	return InviteLinkExpiry(dt), nil
}

// Code is a high-entropy invite link credential: knowing it is what grants
// access, not just a lookup key, so it's generated with crypto/rand rather
// than a predictable id.
type Code string

func NewCode() (Code, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return Code(base64.RawURLEncoding.EncodeToString(b)), nil
}

func (c Code) String() string {
	return string(c)
}

func (d DateTime) String() string {
	return time.Time(d).UTC().Format(time.RFC3339Nano)
}

type UserNotFoundError struct {
	ID ID
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", e.ID)
}

type TeamNotFoundError struct {
	ID ID
}

func (e TeamNotFoundError) Error() string {
	return fmt.Sprintf("team %s not found", e.ID)
}

type NotTeamOwnerError struct {
	ID ID
}

func (e NotTeamOwnerError) Error() string {
	return fmt.Sprintf("caller is not owner of team %s", e.ID)
}

type UserInfo struct {
	ID        ID
	CreatedAt DateTime
}

type TeamInfo struct {
	ID        ID
	Name      TeamName
	CreatedAt DateTime
}

type InviteLinkInfo struct {
	Code      Code
	UseCount  int64
	ExpiresAt *DateTime
	CreatedAt DateTime
}

// User is a reference to a user by id. Constructing it does no I/O and
// cannot fail; whether the id refers to a real user is only known once
// Info or an action method is called.
type User interface {
	Info(ctx context.Context) (UserInfo, error)
	CreateTeam(ctx context.Context, name TeamName) (Team, error)
}

// Team is a reference to a team by id. Constructing it does no I/O and
// cannot fail; whether the id refers to a real team is only known once
// Info is called.
type Team interface {
	Info(ctx context.Context) (TeamInfo, error)
	CreateInviteLink(
		ctx context.Context,
		actor ID,
		expiresAt *InviteLinkExpiry,
	) (InviteLinkInfo, error)
}

// World is the entry point into the domain: it hands out references to
// aggregates by id, without touching the database.
type World interface {
	User(id ID) User
	Team(id ID) Team
}

type SQLWorld struct {
	pool *pgxpool.Pool
}

func NewSQLWorld(pool *pgxpool.Pool) SQLWorld {
	return SQLWorld{pool: pool}
}

func (w SQLWorld) User(id ID) User {
	return SQLUser{pool: w.pool, id: id}
}

func (w SQLWorld) Team(id ID) Team {
	return SQLTeam{pool: w.pool, id: id}
}

type SQLUser struct {
	pool *pgxpool.Pool
	id   ID
}

func (u SQLUser) Info(ctx context.Context) (UserInfo, error) {
	row, err := db.New(u.pool).GetUser(ctx, uuid.UUID(u.id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserInfo{}, UserNotFoundError{ID: u.id}
		}
		return UserInfo{}, fmt.Errorf("get user: %w", err)
	}
	return UserInfo{ID: ID(row.ID), CreatedAt: DateTime(row.CreatedAt)}, nil
}

func (u SQLUser) CreateTeam(ctx context.Context, name TeamName) (Team, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	team, err := db.New(tx).InsertTeam(ctx, string(name))
	if err != nil {
		return nil, fmt.Errorf("insert team: %w", err)
	}

	_, err = db.New(tx).InsertMember(ctx, db.InsertMemberParams{
		TeamID: team.ID,
		UserID: uuid.UUID(u.id),
		Role:   "owner",
	})
	if err != nil {
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok && pgErr.Code == foreignKeyViolation &&
			pgErr.ConstraintName == membersUserIDForeignKey {
			return nil, UserNotFoundError{ID: u.id}
		}
		return nil, fmt.Errorf("insert owner member: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return SQLTeam{
		pool: u.pool,
		id:   ID(team.ID),
		info: &TeamInfo{
			ID:        ID(team.ID),
			Name:      TeamName(team.Name),
			CreatedAt: DateTime(team.CreatedAt),
		},
	}, nil
}

type SQLTeam struct {
	pool *pgxpool.Pool
	id   ID
	// info caches data already known at construction time (e.g. right
	// after an insert), so Info doesn't re-fetch what was just written.
	info *TeamInfo
}

func (t SQLTeam) Info(ctx context.Context) (TeamInfo, error) {
	if t.info != nil {
		return *t.info, nil
	}
	row, err := db.New(t.pool).GetTeam(ctx, uuid.UUID(t.id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TeamInfo{}, TeamNotFoundError{ID: t.id}
		}
		return TeamInfo{}, fmt.Errorf("get team: %w", err)
	}
	return TeamInfo{
		ID:        ID(row.ID),
		Name:      TeamName(row.Name),
		CreatedAt: DateTime(row.CreatedAt),
	}, nil
}

func (t SQLTeam) CreateInviteLink(
	ctx context.Context,
	actor ID,
	expiresAt *InviteLinkExpiry,
) (InviteLinkInfo, error) {
	_, err := t.Info(ctx)
	if err != nil {
		return InviteLinkInfo{}, err
	}

	member, err := db.New(t.pool).GetMember(ctx, db.GetMemberParams{
		TeamID: uuid.UUID(t.id),
		UserID: uuid.UUID(actor),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InviteLinkInfo{}, NotTeamOwnerError{ID: t.id}
		}
		return InviteLinkInfo{}, fmt.Errorf("get member: %w", err)
	}
	if member.Role != "owner" {
		return InviteLinkInfo{}, NotTeamOwnerError{ID: t.id}
	}

	code, err := NewCode()
	if err != nil {
		return InviteLinkInfo{}, fmt.Errorf("generate code: %w", err)
	}

	var expiresAtParam *time.Time
	if expiresAt != nil {
		expiresAtParam = new(time.Time(*expiresAt))
	}

	row, err := db.New(t.pool).InsertInviteLink(ctx, db.InsertInviteLinkParams{
		TeamID:            uuid.UUID(t.id),
		Code:              string(code),
		CreatedByMemberID: member.ID,
		ExpiresAt:         expiresAtParam,
	})
	if err != nil {
		return InviteLinkInfo{}, fmt.Errorf("insert invite link: %w", err)
	}

	var rowExpiresAt *DateTime
	if row.ExpiresAt != nil {
		rowExpiresAt = new(DateTime(*row.ExpiresAt))
	}

	return InviteLinkInfo{
		Code:      Code(row.Code),
		UseCount:  row.UseCount,
		ExpiresAt: rowExpiresAt,
		CreatedAt: DateTime(row.CreatedAt),
	}, nil
}
