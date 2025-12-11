package keyvalue

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/core"
	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
)

type keyValueSession struct {
	key        string
	dataSource keyvalue.Store
}

func (k keyValueSession) StateID(ctx context.Context) (core.ID, error) {
	exists, err := k.dataSource.Exists(ctx, k.stateIDKey())
	if err != nil {
		return core.Stay, fmt.Errorf("checking session existence: %w", err)
	}
	if !exists {
		return core.Stay, nil
	}
	stateID, err := k.dataSource.Read(ctx, k.stateIDKey())
	if err != nil {
		return core.Stay, fmt.Errorf("reading session state ID from %q: %w", k.stateIDKey(), err)
	}
	return core.ID(stateID), nil
}

func (k keyValueSession) ChangeStateID(ctx context.Context, id core.ID) error {
	err := k.dataSource.Write(ctx, k.stateIDKey(), string(id))
	if err != nil {
		return fmt.Errorf("writing session state ID %q to %q: %w", id, k.stateIDKey(), err)
	}
	return nil
}

func (k keyValueSession) Exists(ctx context.Context) (bool, error) {
	exists, err := k.dataSource.Exists(ctx, k.stateIDKey())
	if err != nil {
		return false, fmt.Errorf("checking existence of session state ID at %q: %w", k.stateIDKey(), err)
	}
	return exists, nil
}

func (k keyValueSession) stateIDKey() string {
	return fmt.Sprintf("%s:session_id", k.key)
}
