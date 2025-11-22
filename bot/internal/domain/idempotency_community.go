package domain

import (
	"context"
	"errors"

	"github.com/andreychh/coopera-bot/internal/domain/transport"
)

type idempotencyCommunity struct {
	origin Community
}

func (i idempotencyCommunity) CreateUser(ctx context.Context, tgID int64, tgUsername string) (User, error) {
	_, err := i.origin.CreateUser(ctx, tgID, tgUsername)
	if errors.Is(err, transport.ErrRecordAlreadyExists) {
		return i.origin.UserWithTelegramID(tgID), nil
	}
	return nil, err
}

func (i idempotencyCommunity) UserWithTelegramID(tgID int64) User {
	return i.origin.UserWithTelegramID(tgID)
}

func (i idempotencyCommunity) Team(id int64) Team {
	return i.origin.Team(id)
}

func IdempotencyCommunity(origin Community) Community {
	return idempotencyCommunity{
		origin: origin,
	}
}
