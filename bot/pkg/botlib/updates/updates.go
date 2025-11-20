package updates

type UpdateType int

const (
	UpdateTypeUnknown = iota
	UpdateTypeMessage
	UpdateTypeCallbackQuery
)
