package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
)

func (h *Handler) createTeam(ctx context.Context, input *createTeamInput) (*teamOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	team, err := h.service.CreateTeam(ctx, appteam.CreateTeamInput{
		CurrentUserID: currentUserID,
		Name:          input.Body.Name,
		Slug:          input.Body.Slug,
	})
	if err != nil {
		return nil, mapTeamError(err, "create team failed")
	}

	return &teamOutput{Body: teamOutputBody{Team: newTeamResponse(team)}}, nil
}

type createTeamInput struct {
	Body createTeamInputBody
}

type createTeamInputBody struct {
	Name string `json:"name" minLength:"1" maxLength:"100" required:"true" doc:"Team display name." example:"Engineering"`
	Slug string `json:"slug" minLength:"1" maxLength:"100" pattern:"^[a-z0-9-]+$" patternDescription:"lowercase letters, numbers, and dashes" required:"true" doc:"Stable team slug." example:"engineering"`
}
