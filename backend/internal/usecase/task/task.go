package task

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera-backend/internal/entity"
	"github.com/andreychh/coopera-backend/internal/usecase"
)

type TaskUsecase struct {
	txManager          usecase.TransactionManageRepository
	taskRepository     usecase.TaskRepository
	membershipsUsecase usecase.MembershipUseCase
}

func NewTaskUsecase(taskRepo usecase.TaskRepository, membershipsUsecase usecase.MembershipUseCase, txManager usecase.TransactionManageRepository) *TaskUsecase {
	return &TaskUsecase{
		txManager:          txManager,
		membershipsUsecase: membershipsUsecase,
		taskRepository:     taskRepo,
	}
}

func (uc *TaskUsecase) CreateUsecase(ctx context.Context, task entity.Task) (entity.Task, error) {
	var createdTask entity.Task

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		t, err := uc.taskRepository.CreateRepo(txCtx, task)
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}
		createdTask = t
		return nil
	})

	if err != nil {
		return entity.Task{}, err
	}
	return createdTask, nil
}
