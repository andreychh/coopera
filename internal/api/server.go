// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

import (
	"github.com/andreychh/coopera/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pool *pgxpool.Pool
	post usecase.Post
	seal usecase.Seal
}

func NewServer(pool *pgxpool.Pool, post usecase.Post, seal usecase.Seal) Server {
	return Server{pool: pool, post: post, seal: seal}
}
