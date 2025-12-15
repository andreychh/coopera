package domain

type UserStats struct {
	Teams map[string]UserTeamStats
}

type UserTeamStats struct {
	ActiveLoad struct {
		TasksCount  int
		TotalPoints int
	}

	LifetimeContribution struct {
		TasksCompleted int
		PointsEarned   int
	}
}

type MemberStats struct {
	CurrentState struct {
		AssignedTasks  int
		AssignedPoints int
	}
	Contribution struct {
		TasksDone  int
		PointsDone int
	}
}

type TeamStats struct {
	Backlog struct {
		DraftsCount      int
		UnassignedCount  int
		UnassignedPoints int
	}
	ActiveWork struct {
		InProgressCount  int
		InProgressPoints int
		InReviewCount    int
		InReviewPoints   int
	}
	Achievements struct {
		CompletedCount  int
		CompletedPoints int
	}
}
