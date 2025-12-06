package memory

import "github.com/andreychh/coopera-bot/internal/domain"

type userDetails struct {
	id int64
}

func (u userDetails) ID() int64 {
	return u.id
}

func UserDetails(id int64) domain.UserDetails {
	return userDetails{
		id: id,
	}
}
