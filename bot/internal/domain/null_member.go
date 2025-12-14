package domain

import "context"

type nullMember struct{}

func (n nullMember) ID() int64 {
	panic("not implemented")
}

func (n nullMember) Name() string {
	panic("not implemented")
}

func (n nullMember) Role() MemberRole {
	panic("not implemented")
}

func (n nullMember) CreateTask(ctx context.Context, title string, description string, points int, assignee Member) (Task, error) {
	panic("not implemented")
}

func (n nullMember) Tasks(ctx context.Context) (Tasks, error) {
	panic("not implemented")
}

func NullMember() Member {
	return nullMember{}
}
