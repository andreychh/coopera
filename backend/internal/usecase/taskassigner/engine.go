package taskassigner

import (
	"context"
	"errors"

	"github.com/andreychh/coopera-backend/internal/entity"
	"github.com/andreychh/coopera-backend/internal/usecase"
)

type taskAssignmentUsecase struct {
	taskUsecase    usecase.TaskUseCase
	membersUsecase usecase.MembershipUseCase
	txRepo         usecase.TransactionManageRepository
}

func NewTaskAssignmentUsecase(
	taskUsecase usecase.TaskUseCase,
	membersUsecase usecase.MembershipUseCase,
	txRepo usecase.TransactionManageRepository,
) usecase.TaskAssignmentUsecase {
	return &taskAssignmentUsecase{
		taskUsecase:    taskUsecase,
		membersUsecase: membersUsecase,
		txRepo:         txRepo,
	}
}

func (u *taskAssignmentUsecase) AssignTasks(ctx context.Context) error {
	return u.txRepo.WithinTransaction(ctx, func(txCtx context.Context) error {
		tasks, err := u.taskUsecase.GetAllTasks(txCtx)
		if err != nil {
			return err
		}

		// валидация по статусу open - сутки (для презы поставить 20 секунд)
		// TODO: в будущем подумать над тем что у задачи может не быть баллов

		for _, task := range tasks {
			if task.AssignedToMember != nil {
				continue
			}

			users, err := u.membersUsecase.GetMembersUsecase(txCtx, task.TeamID)
			if err != nil {
				return err
			}

			if len(users) == 0 {
				return errors.New("no users available for assignment")
			}

			var targetUser entity.MembershipEntity
			minLoad := int32(1<<31 - 1)
			for _, user := range users {
				userTasks, err := u.taskUsecase.GetUsecase(txCtx, entity.TaskFilter{UserID: user.ID})
				if err != nil {
					return err
				}

				var totalPoints int32
				for _, t := range userTasks {
					if *t.Status == entity.StatusAssigned || *t.Status == entity.StatusInReview {
						totalPoints += *t.Points
					}
				}

				load := totalPoints + int32(len(userTasks))
				if load < minLoad {
					minLoad = load
					targetUser = user
				}
			}

			if err := u.taskUsecase.UpdateForEngine(txCtx, entity.UpdateTask{
				TaskID:           task.ID,
				AssignedToMember: &targetUser.UserID,
				Description:      task.Description,
			}); err != nil {
				return err
			}
			if err := u.taskUsecase.UpdateStatusForEngine(txCtx, entity.TaskStatus{
				TaskID: task.ID,
				Status: entity.StatusAssigned.String(),
			}); err != nil {
				return err
			}
		}

		return nil
	})
}
