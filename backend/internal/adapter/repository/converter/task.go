package converter

import (
	taskModel "github.com/andreychh/coopera/internal/adapter/repository/model/task_model"
	"github.com/andreychh/coopera/internal/entity"
)

func FromEntityToModelTask(task entity.Task) taskModel.Task {
	mtask := taskModel.Task{
		ID:          task.ID,
		TeamID:      task.TeamID,
		Title:       task.Title,
		Description: task.Description,
		Points:      task.Points,
		AssignedTo:  task.AssignedTo,
		CreatedBy:   task.CreatedBy,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if task.Status != nil {
		mtask.Status = task.Status.String()
	} else {
		mtask.Status = entity.StatusOpen.String()
	}

	return mtask
}

func FromModelToEntityTask(m taskModel.Task) entity.Task {
	status := entity.Status(m.Status)
	return entity.Task{
		ID:          m.ID,
		TeamID:      m.TeamID,
		Title:       m.Title,
		Description: m.Description,
		Points:      m.Points,
		Status:      &status,
		AssignedTo:  m.AssignedTo,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
