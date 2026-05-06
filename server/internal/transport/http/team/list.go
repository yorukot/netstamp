package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

func (h *Handler) listTeams(ctx context.Context, _ *listTeamsInput) (*listTeamsOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	teams, err := h.service.ListTeams(ctx, appteam.ListTeamsInput{CurrentUserID: currentUserID})
	if err != nil {
		return nil, mapTeamError(err, "list teams failed")
	}

	return &listTeamsOutput{Body: listTeamsOutputBody{Teams: newTeamResponses(teams)}}, nil
}

type listTeamsInput struct{}

type listTeamsOutput struct {
	Body listTeamsOutputBody
}

type listTeamsOutputBody struct {
	Teams []teamResponse `json:"teams"`
}

func newTeamResponses(teams []domainteam.Team) []teamResponse {
	responses := make([]teamResponse, 0, len(teams))
	for _, team := range teams {
		responses = append(responses, newTeamResponse(team))
	}

	return responses
}
