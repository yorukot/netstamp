package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
)

func (h *Handler) updateTeam(ctx context.Context, input *updateTeamInput) (*teamOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	team, err := h.service.UpdateTeam(ctx, appteam.UpdateTeamInput{
		CurrentUserID: currentUserID,
		TeamRef:       input.Ref,
		Name:          input.Body.Name,
		Slug:          input.Body.Slug,
	})
	if err != nil {
		return nil, mapTeamError(err, "update team failed")
	}

	return &teamOutput{Body: teamOutputBody{Team: newTeamResponse(team)}}, nil
}

type updateTeamInput struct {
	Ref  string `path:"ref" minLength:"1" maxLength:"100" doc:"Team UUID or slug." example:"engineering"`
	Body updateTeamInputBody
}

type updateTeamInputBody struct {
	Name *string `json:"name,omitempty" minLength:"1" maxLength:"100" doc:"Team display name." example:"Engineering"`
	Slug *string `json:"slug,omitempty" minLength:"1" maxLength:"100" pattern:"^[a-z0-9-]+$" patternDescription:"lowercase letters, numbers, and dashes" doc:"Stable team slug." example:"engineering"`
}
