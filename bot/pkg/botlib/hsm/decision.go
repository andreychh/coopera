package hsm

import "github.com/andreychh/coopera-bot/pkg/botlib/core"

type Decision interface {
	Next() core.ID
	Handled() bool
	HasTransition() bool
}

type decision struct {
	next    core.ID
	handled bool
}

func (d decision) Next() core.ID {
	return d.next
}

func (d decision) Handled() bool {
	return d.handled
}

func (d decision) HasTransition() bool {
	return d.handled && d.next != core.Stay
}

func Pass() Decision {
	return decision{
		next:    core.Stay,
		handled: false,
	}
}

func Transit(next core.ID) Decision {
	return decision{
		next:    next,
		handled: true,
	}
}

func Stay() Decision {
	return decision{
		next:    core.Stay,
		handled: true,
	}
}
