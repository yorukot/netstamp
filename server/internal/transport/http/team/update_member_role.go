package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

func (h *Handler) updateMemberRole(ctx context.Context, input *updateMemberRoleInput) (*memberOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	member, err := h.service.UpdateMemberRole(ctx, appteam.UpdateMemberRoleInput{
		CurrentUserID: currentUserID,
		TeamID:        input.ID,
		UserID:        input.UserID,
		Role:          domainteam.Role(input.Body.Role),
	})
	if err != nil {
		return nil, mapTeamError(err, "update team member failed")
	}

	return &memberOutput{Body: memberOutputBody{Member: newTeamMemberResponse(member)}}, nil
}

type updateMemberRoleInput struct {
	ID     string `path:"id" format:"uuid"`
	UserID string `path:"user_id" format:"uuid"`
	Body   updateMemberRoleInputBody
}

type updateMemberRoleInputBody struct {
	Role string `json:"role" enum:"owner,admin,editor,viewer" required:"true" doc:"Team member role." example:"viewer"`
}
