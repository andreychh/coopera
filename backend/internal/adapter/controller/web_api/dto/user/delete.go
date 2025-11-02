package user

type DeleteUserRequest struct {
	TelegramID int64 `json:"telegram_id" validate:"required"`
}

func (r *DeleteUserRequest) GetTelegramID() int64 { return r.TelegramID }
