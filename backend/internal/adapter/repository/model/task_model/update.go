package task_model

type UpdateTask struct {
	ID          int32   `db:"id"`
	Title       *string `db:"title"`
	Description *string `db:"description"`
	Points      *int32  `db:"points"`
	AssignedTo  *int32  `db:"assigned_to"`
}
