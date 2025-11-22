package updates

type UpdateType string

const (
	UpdateTypeUnknown       UpdateType = "unknown"
	UpdateTypeMessage       UpdateType = "message"
	UpdateTypeCallbackQuery UpdateType = "callback_query"
)

type ChatType string

const (
	ChatTypeUnknown    ChatType = "unknown"
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeSuperGroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)
