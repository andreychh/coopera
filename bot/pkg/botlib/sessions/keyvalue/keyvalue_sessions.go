package keyvalue

import (
	"fmt"

	"github.com/andreychh/coopera-bot/pkg/botlib/keyvalue"
	"github.com/andreychh/coopera-bot/pkg/botlib/sessions"
)

type keyValueSessions struct {
	dataSource keyvalue.Store
}

func (k keyValueSessions) Session(id int64) sessions.Session {
	return keyValueSession{
		key:        k.sessionKey(id),
		dataSource: k.dataSource,
	}
}

func (k keyValueSessions) sessionKey(id int64) string {
	return fmt.Sprintf("session:%d", id)
}

func Sessions(dataSource keyvalue.Store) sessions.Sessions {
	return keyValueSessions{dataSource: dataSource}
}
