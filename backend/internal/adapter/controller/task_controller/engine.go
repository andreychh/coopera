package task_controller

import (
	"context"
	"log"
	"time"

	"github.com/andreychh/coopera-backend/internal/usecase"
)

type TaskAssignmentController struct {
	usecase usecase.TaskAssignmentUsecase
}

func NewTaskAssignmentController(u usecase.TaskAssignmentUsecase) *TaskAssignmentController {
	return &TaskAssignmentController{
		usecase: u,
	}
}

func (c *TaskAssignmentController) StartAssignmentLoop(ctx context.Context, interval time.Duration, taskMinAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("task assignment loop stopped")
			return
		case <-ticker.C:
			loopCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
			if err := c.usecase.AssignTasks(loopCtx, taskMinAge); err != nil {
				log.Println("failed to assign tasks:", err)
			}
			cancel()
		}
	}
}
