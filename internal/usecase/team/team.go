package team

import (
	"context"
	"github.com/andreychh/coopera/internal/entity"
	"github.com/andreychh/coopera/internal/usecase"
)

type TeamUsecase struct {
	txManager          usecase.TransactionManager
	teamRepository     usecase.TeamRepository
	membershipsUsecase usecase.MembershipUseCase
}

func NewTeamUsecase(teamRepo usecase.TeamRepository, membershipsUsecase usecase.MembershipUseCase, txManager usecase.TransactionManager) *TeamUsecase {
	return &TeamUsecase{
		txManager:          txManager,
		membershipsUsecase: membershipsUsecase,
		teamRepository:     teamRepo,
	}
}

func (uc *TeamUsecase) CreateUsecase(ctx context.Context, team entity.TeamEntity) (entity.TeamEntity, error) {
	var createdTeam entity.TeamEntity

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		t, err := uc.teamRepository.CreateRepo(txCtx, team)
		if err != nil {
			return err
		}
		createdTeam = t

		return uc.membershipsUsecase.AddMemberUsecase(txCtx, entity.MembershipEntity{
			TeamID:   *t.ID,
			MemberID: team.CreatedBy,
			Role:     entity.RoleManager,
		})
	})

	if err != nil {
		return entity.TeamEntity{}, err
	}
	return createdTeam, nil
}

// Получение команды и участников
func (uc *TeamUsecase) GetByIDUsecase(ctx context.Context, teamID int32) (entity.TeamEntity, []entity.MembershipEntity, error) {
	var (
		team    entity.TeamEntity
		members []entity.MembershipEntity
	)

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		team, err = uc.teamRepository.GetByIDRepo(txCtx, teamID)
		if err != nil {
			return err
		}

		members, err = uc.membershipsUsecase.GetMembersUsecase(txCtx, teamID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.TeamEntity{}, nil, err
	}

	return team, members, nil
}
