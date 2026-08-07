// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import "github.com/jackc/pgx/v5/pgxpool"

type Server struct {
	pool *pgxpool.Pool
}

func NewServer(pool *pgxpool.Pool) Server {
	return Server{pool: pool}
}
