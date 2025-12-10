package domain

import (
	"context"

	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
	"github.com/andreychh/coopera-bot/pkg/botlib/updates/attributes"
)

func UserWithID(community Community, id sources.Source[int64]) sources.Source[User] {
	return sources.PureMap(id, community.UserWithTelegramID)
}

func CurrentUser(community Community) sources.Source[User] {
	return UserWithID(community, sources.Required(attributes.ChatID()))
}

func CurrentTeams(community Community) sources.Source[Teams] {
	return sources.PureMap(CurrentUser(community),
		func(user User) Teams {
			return user.CreatedTeams()
		},
	)
}

func TeamWithName(teams sources.Source[Teams], name sources.Source[string]) sources.Source[Team] {
	return sources.PureZip(teams, name,
		func(teams Teams, name string) Team {
			return teams.TeamWithName(name)
		},
	)
}

func DetailsOfTeam(team sources.Source[Team]) sources.Source[TeamDetails] {
	return sources.Map(team,
		func(ctx context.Context, team Team) (TeamDetails, error) {
			return team.Details(ctx)
		},
	)
}
