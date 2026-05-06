package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
)

func (h *Handler) getTeam(ctx context.Context, input *teamIDInput) (*teamOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	team, err := h.service.GetTeam(ctx, appteam.GetTeamInput{
		CurrentUserID: currentUserID,
		TeamID:        input.ID,
	})
	if err != nil {
		return nil, mapTeamError(err, "get team failed")
	}

	return &teamOutput{Body: teamOutputBody{Team: newTeamResponse(team)}}, nil
}

type teamIDInput struct {
	ID string `path:"id" format:"uuid"`
}
