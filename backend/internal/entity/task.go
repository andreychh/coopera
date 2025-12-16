package entity

import "time"

type Task struct {
	ID                int32
	TeamID            int32
	Title             string
	Description       *string
	Points            *int32
	Status            *Status
	AssignedToMember  *int32
	CreatedByMemberID int32
	CreatedByUserID   int32
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}
type TaskFilter struct {
	TaskID   int32
	MemberID int32
	TeamID   int32
}

type TaskStatus struct {
	TaskID        int32
	Status        string
	CurrentUserID int32
}

type UpdateTask struct {
	TaskID           int32
	Title            *string
	Description      *string
	Points           *int32
	AssignedToMember *int32
}
