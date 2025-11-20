package entity

type Role string

const (
	RoleManager Role = "manager"
	RoleMember  Role = "member"
)

func (r Role) IsValid() bool {
	return r == RoleManager || r == RoleMember
}

type Status string

const (
	StatusOpen      Status = "open"
	StatusAssigned  Status = "assigned"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

// String возвращает строковое представление статуса
func (s Status) String() string {
	return string(s)
}
