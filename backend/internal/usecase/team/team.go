package team

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera-backend/internal/entity"
	"github.com/andreychh/coopera-backend/internal/usecase"
	appErr "github.com/andreychh/coopera-backend/pkg/errors"
	"github.com/pkg/errors"
)

type TeamUsecase struct {
	txManager          usecase.TransactionManageRepository
	teamRepository     usecase.TeamRepository
	membershipsUsecase usecase.MembershipUseCase
	userUsecase        usecase.UserUseCase
}

func NewTeamUsecase(teamRepo usecase.TeamRepository, membershipsUsecase usecase.MembershipUseCase, userUsecase usecase.UserUseCase, txManager usecase.TransactionManageRepository) *TeamUsecase {
	return &TeamUsecase{
		txManager:          txManager,
		membershipsUsecase: membershipsUsecase,
		teamRepository:     teamRepo,
		userUsecase:        userUsecase,
	}
}

func (uc *TeamUsecase) CreateUsecase(ctx context.Context, team entity.TeamEntity) (entity.TeamEntity, error) {
	var createdTeam entity.TeamEntity

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		exists, err := uc.teamRepository.ExistsByName(txCtx, team.Name)
		if err != nil {
			return err
		}
		if exists {
			return errors.Wrap(appErr.ErrAlreadyExists, fmt.Sprintf("team '%s'", team.Name))
		}

		t, err := uc.teamRepository.CreateRepo(txCtx, team)
		if err != nil {
			return fmt.Errorf("failed to create team: %w", err)
		}
		createdTeam = t

		_, err = uc.membershipsUsecase.AddMemberUsecase(txCtx, entity.MembershipEntity{
			TeamID: *t.ID,
			UserID: team.CreatedBy,
			Role:   entity.RoleManager,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return entity.TeamEntity{}, err
	}
	return createdTeam, nil
}

func (uc *TeamUsecase) DeleteUsecase(ctx context.Context, teamID, currentUserID int32) error {
	return uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		members, err := uc.membershipsUsecase.GetMembersUsecase(txCtx, teamID)
		if err != nil {
			return err
		}

		var currentUserRole entity.Role
		for _, m := range members {
			if m.UserID == currentUserID {
				currentUserRole = m.Role
				break
			}
		}

		if currentUserRole != entity.RoleManager {
			return appErr.ErrNoPermissionToDelete
		}

		if err := uc.teamRepository.DeleteRepo(txCtx, teamID); err != nil {
			return fmt.Errorf("failed to delete team: %w", err)
		}
		return nil
	})
}

func (uc *TeamUsecase) GetByIDUsecase(ctx context.Context, teamID int32) (entity.TeamEntity, []entity.MembershipEntity, map[int32]string, error) {

	var (
		team      entity.TeamEntity
		members   []entity.MembershipEntity
		usernames = make(map[int32]string)
	)

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error

		team, err = uc.teamRepository.GetByIDRepo(txCtx, teamID)
		if err != nil {
			return fmt.Errorf("failed to get team: %w", err)
		}

		members, err = uc.membershipsUsecase.GetMembersUsecase(txCtx, teamID)
		if err != nil {
			return err
		}

		uniqUserIDs := make(map[int32]struct{})
		for _, m := range members {
			uniqUserIDs[m.UserID] = struct{}{}
		}
		uniqUserIDs[team.CreatedBy] = struct{}{}

		for userID := range uniqUserIDs {
			user, err := uc.userUsecase.GetUsecase(txCtx, userID)
			if err != nil {
				return err
			}
			usernames[userID] = *user.Username
		}

		return nil
	})

	if err != nil {
		return entity.TeamEntity{}, nil, nil, err
	}

	return team, members, usernames, nil
}

func (uc *TeamUsecase) ExistTeamByIDUsecase(ctx context.Context, teamID int32) (bool, error) {
	var exists bool
	var err error

	err = uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		exists, err = uc.teamRepository.ExistsByID(txCtx, teamID)
		if err != nil {
			return fmt.Errorf("failed to check team existence: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}
