package entity

import "time"

type MembershipEntity struct {
	ID        int32
	TeamID    int32
	MemberID  int32
	Role      Role
	CreatedAt *time.Time
}
