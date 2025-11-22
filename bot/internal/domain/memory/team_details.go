package memory

import "github.com/andreychh/coopera-bot/internal/domain"

type teamDetails struct {
	id   int64
	name string
}

func (t teamDetails) ID() int64 {
	return t.id
}

func (t teamDetails) Name() string {
	return t.name
}

func TeamDetails(id int64, name string) domain.TeamDetails {
	return teamDetails{
		id:   id,
		name: name,
	}
}
