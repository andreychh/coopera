package task

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera-backend/internal/entity"
	"github.com/andreychh/coopera-backend/internal/usecase"
	appErr "github.com/andreychh/coopera-backend/pkg/errors"
)

type TaskUsecase struct {
	txManager          usecase.TransactionManageRepository
	taskRepository     usecase.TaskRepository
	membershipsUsecase usecase.MembershipUseCase
	teamUsecase        usecase.TeamUseCase
}

func NewTaskUsecase(taskRepo usecase.TaskRepository, membershipsUsecase usecase.MembershipUseCase, txManager usecase.TransactionManageRepository, teamUsecase usecase.TeamUseCase) *TaskUsecase {
	return &TaskUsecase{
		txManager:          txManager,
		membershipsUsecase: membershipsUsecase,
		taskRepository:     taskRepo,
		teamUsecase:        teamUsecase,
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

func (uc *TaskUsecase) GetUsecase(ctx context.Context, f entity.TaskFilter) ([]entity.Task, error) {
	var result []entity.Task

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {

		switch {
		case f.TaskID > 0:
			task, err := uc.taskRepository.GetByTaskID(txCtx, f.TaskID)
			if err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}
			result = []entity.Task{task}
			return nil

		case f.UserID > 0:
			exists, err := uc.membershipsUsecase.ExistsMemberUsecase(ctx, f.UserID)
			if err != nil {
				return fmt.Errorf("failed to check membership: %w", err)
			}
			if exists {
				tasks, err := uc.taskRepository.GetByAssignedToID(txCtx, f.UserID)
				if err != nil {
					return fmt.Errorf("failed to get task: %w", err)
				}
				result = tasks
				return nil
			}

			return appErr.ErrMemberNotFound

		case f.TeamID > 0:
			exists, err := uc.teamUsecase.ExistTeamByIDUsecase(ctx, f.TeamID)
			if err != nil {
				return fmt.Errorf("failed to check team: %w", err)
			}
			if exists {
				tasks, err := uc.taskRepository.GetByTeamID(txCtx, f.TeamID)
				if err != nil {
					return fmt.Errorf("failed to get task: %w", err)
				}
				result = tasks
				return nil
			}

			return appErr.ErrTeamNotFound

		default:
			return appErr.ErrTaskFilter
		}
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
