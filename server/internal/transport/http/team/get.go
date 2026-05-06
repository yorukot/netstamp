package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
)

func (h *Handler) getTeam(ctx context.Context, input *teamRefInput) (*teamOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	team, err := h.service.GetTeam(ctx, appteam.GetTeamInput{
		CurrentUserID: currentUserID,
		TeamRef:       input.Ref,
	})
	if err != nil {
		return nil, mapTeamError(err, "get team failed")
	}

	return &teamOutput{Body: teamOutputBody{Team: newTeamResponse(team)}}, nil
}

type teamRefInput struct {
	Ref string `path:"ref" minLength:"1" maxLength:"100" doc:"Team UUID or slug." example:"engineering"`
}
