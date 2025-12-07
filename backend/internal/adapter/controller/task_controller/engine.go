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

func (c *TaskAssignmentController) StartAssignmentLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.usecase.AssignTasks(ctx); err != nil {
				log.Println("failed to assign tasks:", err)
			}
		}
	}
}
