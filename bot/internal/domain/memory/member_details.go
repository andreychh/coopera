package memory

import "github.com/andreychh/coopera-bot/internal/domain"

type memberDetails struct {
	id   int64
	name string
}

func (m memberDetails) ID() int64 {
	return m.id
}

func (m memberDetails) Name() string {
	return m.name
}

func MemberDetails(id int64, name string) domain.MemberDetails {
	return memberDetails{
		id:   id,
		name: name,
	}
}
