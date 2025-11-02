package entity

type Role string

const (
	RoleManager Role = "manager"
	RoleMember  Role = "member"
)

func (r Role) IsValid() bool {
	return r == RoleManager || r == RoleMember
}
