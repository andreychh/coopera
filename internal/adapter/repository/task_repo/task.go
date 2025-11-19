package task_repo

import (
	"context"
	"github.com/andreychh/coopera/internal/adapter/repository/converter"
	"github.com/andreychh/coopera/internal/adapter/repository/postgres/dao"
	"github.com/andreychh/coopera/internal/entity"
)

type TaskRepository struct {
	TaskDAO dao.TaskDAO
}

func NewTaskRepository(taskDAO dao.TaskDAO) *TaskRepository {
	return &TaskRepository{
		TaskDAO: taskDAO,
	}
}

func (ur *TaskRepository) CreateRepo(ctx context.Context, task entity.Task) (entity.Task, error) {
	taskModel := converter.FromEntityToModelTask(task)
	entask, err := ur.TaskDAO.Create(ctx, taskModel)
	if err != nil {
		return entity.Task{}, err
	}

	return entask, nil
}
