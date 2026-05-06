package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

func (h *Handler) addMember(ctx context.Context, input *addMemberInput) (*memberOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	member, err := h.service.AddMember(ctx, appteam.AddMemberInput{
		CurrentUserID: currentUserID,
		TeamRef:       input.Ref,
		UserID:        input.Body.UserID,
		Role:          domainteam.Role(input.Body.Role),
	})
	if err != nil {
		return nil, mapTeamError(err, "add team member failed")
	}

	return &memberOutput{Body: memberOutputBody{Member: newTeamMemberResponse(member)}}, nil
}

type addMemberInput struct {
	Ref  string `path:"ref" minLength:"1" maxLength:"100" doc:"Team UUID or slug." example:"engineering"`
	Body addMemberInputBody
}

type addMemberInputBody struct {
	UserID string `json:"userId" format:"uuid" required:"true" doc:"User ID to add to the team."`
	Role   string `json:"role" enum:"owner,admin,editor,viewer" required:"true" doc:"Team member role." example:"viewer"`
}
