package task_repo

import (
	"context"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/converter"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/postgres/dao"
	"github.com/andreychh/coopera-backend/internal/entity"
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

func (ur *TaskRepository) GetByTaskID(ctx context.Context, id int32) (entity.Task, error) {
	return ur.TaskDAO.GetByTaskID(ctx, id)
}

func (ur *TaskRepository) GetByAssignedToID(ctx context.Context, userID int32) ([]entity.Task, error) {
	return ur.TaskDAO.GetByAssignedID(ctx, userID)
}

func (ur *TaskRepository) GetByTeamID(ctx context.Context, teamID int32) ([]entity.Task, error) {
	return ur.TaskDAO.GetByTeamID(ctx, teamID)
}
