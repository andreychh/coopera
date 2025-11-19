package task_model

import "time"

type Task struct {
	ID          int32      `db:"id"`
	TeamID      int32      `db:"team_id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	Points      int32      `db:"points"`
	Status      string     `db:"status"`
	AssignedTo  *int32     `db:"assigned_to"`
	CreatedBy   int32      `db:"created_by"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
