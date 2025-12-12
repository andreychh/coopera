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
		currentMember, err := uc.membershipsUsecase.GetMemberUsecase(txCtx, task.TeamID, task.CreatedBy)
		if err != nil {
			return fmt.Errorf("failed to get current member: %w", err)
		}

		isManager := currentMember.Role == entity.RoleManager
		if !isManager && task.Points != nil {
			return appErr.ErrOnlyManagerCanSetPoints
		}

		if task.AssignedToMember != nil {
			if !isManager {
				return appErr.ErrOnlyManagerCanAssign
			}

			exists, err := uc.membershipsUsecase.ExistsMemberUsecase(txCtx, *task.AssignedToMember)
			if err != nil {
				return err
			}
			if !exists {
				return appErr.ErrAssignedMemberNotExists
			}
		}

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

func (uc *TaskUsecase) UpdateUsecase(ctx context.Context, task entity.UpdateTask, currUserID int32) error {
	return uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {

		existingTask, err := uc.taskRepository.GetByTaskID(txCtx, task.TaskID)
		if err != nil {
			return fmt.Errorf("failed to get task: %w", err)
		}

		currentMember, err := uc.membershipsUsecase.GetMemberUsecase(txCtx, existingTask.TeamID, currUserID)
		if err != nil {
			return fmt.Errorf("failed to get current member: %w", err)
		}

		isCreator := existingTask.CreatedBy == currentMember.ID
		isManager := currentMember.Role == entity.RoleManager

		if task.Title != nil || task.Description != nil {
			if !isCreator && !isManager {
				return appErr.ErrNoPermissionToUpdate
			}
		}

		if task.Points != nil && !isManager {
			return appErr.ErrOnlyManagerCanUpdatePoints
		}

		needStatusUpdate := false

		if task.AssignedToMember != nil {

			if existingTask.Points == nil || *existingTask.Points <= 0 {
				return appErr.ErrCantAssignWithoutPoints
			}

			if !isManager {
				if *task.AssignedToMember != currentMember.ID {
					return appErr.ErrOnlyManagerOrSelfCanAssign
				}
			}

			needStatusUpdate = true
		}

		if err = uc.taskRepository.UpdateRepo(txCtx, task); err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		if needStatusUpdate {
			if err = uc.UpdateStatus(txCtx, entity.TaskStatus{
				TaskID:        task.TaskID,
				Status:        entity.StatusAssigned.String(),
				CurrentUserID: currUserID,
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (uc *TaskUsecase) UpdateStatusForEngine(ctx context.Context, taskStatus entity.TaskStatus) error {
	return uc.taskRepository.UpdateStatus(ctx, taskStatus)
}

func (uc *TaskUsecase) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	var tasks []entity.Task

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		t, err := uc.taskRepository.GetAllTasks(txCtx)
		if err != nil {
			return fmt.Errorf("failed to get all tasks: %w", err)
		}
		tasks = t
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (uc *TaskUsecase) UpdateForEngine(ctx context.Context, task entity.UpdateTask) error {
	return uc.taskRepository.UpdateRepo(ctx, task)
}

func (uc *TaskUsecase) UpdateStatus(ctx context.Context, task entity.TaskStatus) error {
	return uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		existingTask, err := uc.taskRepository.GetByTaskID(txCtx, task.TaskID)
		if err != nil {
			return fmt.Errorf("failed to get task: %w", err)
		}

		members, err := uc.membershipsUsecase.GetMembersUsecase(txCtx, existingTask.TeamID)
		if err != nil {
			return err
		}

		isMember := false
		for _, m := range members {
			if m.UserID == task.CurrentUserID {
				isMember = true
				break
			}
		}

		if !isMember {
			return appErr.ErrNoPermissionToUpdate
		}

		err = uc.taskRepository.UpdateStatus(txCtx, task)
		if err != nil {
			return fmt.Errorf("failed to update task status: %w", err)
		}
		return nil
	})
}

func (uc *TaskUsecase) DeleteUsecase(ctx context.Context, taskID, currentUserID int32) error {
	return uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		existingTask, err := uc.taskRepository.GetByTaskID(txCtx, taskID)
		if err != nil {
			return fmt.Errorf("failed to get task: %w", err)
		}

		currentMember, err := uc.membershipsUsecase.GetMemberUsecase(txCtx, existingTask.TeamID, currentUserID)
		if err != nil {
			return fmt.Errorf("failed to get current member: %w", err)
		}

		isCreator := existingTask.CreatedBy == currentMember.ID
		isManager := currentMember.Role == entity.RoleManager

		if !isCreator && !isManager {
			return appErr.ErrNoPermissionToDelete
		}

		err = uc.taskRepository.DeleteRepo(txCtx, taskID)
		if err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		return nil
	})
}
