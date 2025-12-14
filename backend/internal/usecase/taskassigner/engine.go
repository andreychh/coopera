package taskassigner

import (
	"context"
	"github.com/andreychh/coopera-backend/internal/entity"
	"github.com/andreychh/coopera-backend/internal/usecase"
	"time"
)

type taskAssignmentUsecase struct {
	txManager          usecase.TransactionManageRepository
	taskUsecase        usecase.TaskUseCase
	membershipsUsecase usecase.MembershipUseCase
}

func NewTaskAssignmentUsecase(
	txManager usecase.TransactionManageRepository,
	taskUsecase usecase.TaskUseCase,
	membershipsUsecase usecase.MembershipUseCase,
) usecase.TaskAssignmentUsecase {
	return &taskAssignmentUsecase{
		txManager:          txManager,
		taskUsecase:        taskUsecase,
		membershipsUsecase: membershipsUsecase,
	}
}

func (u *taskAssignmentUsecase) AssignTasks(ctx context.Context, taskMinAge time.Duration) error {
	return u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {

		tasks, err := u.taskUsecase.GetAllTasks(txCtx)
		if err != nil {
			return err
		}

		mapTaskTeam := make(map[int32][]entity.Task)
		for _, task := range tasks {
			mapTaskTeam[task.TeamID] = append(mapTaskTeam[task.TeamID], task)
		}

		now := time.Now()

		for teamID, teamTasks := range mapTaskTeam {

			members, err := u.membershipsUsecase.GetMembersUsecase(txCtx, teamID)
			if err != nil {
				return err
			}

			mapUserIDPoints := make(map[int32]int32, len(members))
			for _, m := range members {
				mapUserIDPoints[m.ID] = 0
			}

			var notAssignedTasks []entity.Task
			for _, teamTask := range teamTasks {
				if teamTask.AssignedToMember == nil && now.Sub(*teamTask.CreatedAt) >= taskMinAge && teamTask.Points != nil {
					notAssignedTasks = append(notAssignedTasks, teamTask)
				} else if teamTask.AssignedToMember != nil &&
					(*teamTask.Status == entity.StatusAssigned || *teamTask.Status == entity.StatusInReview) {
					mapUserIDPoints[*teamTask.AssignedToMember] += *teamTask.Points
				}
			}

			for _, task := range notAssignedTasks {
				var minPoints int32 = -1
				var selectedMemberID int32

				for memberID, points := range mapUserIDPoints {
					if minPoints == -1 || points < minPoints {
						minPoints = points
						selectedMemberID = memberID
					}
				}

				if err := u.taskUsecase.UpdateForEngine(txCtx, entity.UpdateTask{
					TaskID:           task.ID,
					AssignedToMember: &selectedMemberID,
				}); err != nil {
					return err
				}

				if err := u.taskUsecase.UpdateStatusForEngine(txCtx, entity.TaskStatus{
					TaskID: task.ID,
					Status: entity.StatusAssigned.String(),
				}); err != nil {
					return err
				}

				mapUserIDPoints[selectedMemberID] += *task.Points
			}
		}

		return nil
	})
}
