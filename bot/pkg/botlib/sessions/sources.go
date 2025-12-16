package sessions

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
)

func SessionWithID(sessions Sessions, id sources.Source[int64]) sources.Source[Session] {
	return sources.PureMap(id,
		func(id int64) Session {
			return sessions.Session(id)
		},
	)
}

func CurrentSession(sessions Sessions) sources.Source[Session] {
	return SessionWithID(sessions, sources.Required(attributes.ChatID()))
}
