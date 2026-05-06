package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
)

func (h *Handler) deleteTeam(ctx context.Context, input *teamIDInput) (*deleteTeamOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.service.DeleteTeam(ctx, appteam.DeleteTeamInput{
		CurrentUserID: currentUserID,
		TeamID:        input.ID,
	}); err != nil {
		return nil, mapTeamError(err, "delete team failed")
	}

	return &deleteTeamOutput{}, nil
}

type deleteTeamOutput struct{}
