package entity

import "time"

type MembershipEntity struct {
	ID        int32
	TeamID    int32
	UserID    int32
	Role      Role
	CreatedAt *time.Time
}
