package task_model

type TaskStatus struct {
	TaskID int32  `db:"task_id"`
	Status string `db:"status"`
}
