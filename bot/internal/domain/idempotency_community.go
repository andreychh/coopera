package domain

import (
	"context"
	"errors"

	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type idempotencyCommunity struct {
	origin Community
}

func (i idempotencyCommunity) CreateUser(ctx context.Context, telegramID int64) (User, error) {
	_, err := i.origin.CreateUser(ctx, telegramID)
	if errors.Is(err, transport.ErrRecordAlreadyExists) {
		return i.origin.UserWithTelegramID(telegramID), nil
	}
	return nil, err
}

func (i idempotencyCommunity) UserWithTelegramID(telegramID int64) User {
	return i.origin.UserWithTelegramID(telegramID)
}

func (i idempotencyCommunity) Team(id int64) Team {
	return i.origin.Team(id)
}

func IdempotencyCommunity(origin Community) Community {
	return idempotencyCommunity{
		origin: origin,
	}
}
