package team

import (
	"context"

	appteam "github.com/yorukot/netstamp/internal/application/team"
	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

func (h *Handler) listMembers(ctx context.Context, input *teamRefInput) (*listMembersOutput, error) {
	currentUserID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	members, err := h.service.ListMembers(ctx, appteam.ListMembersInput{
		CurrentUserID: currentUserID,
		TeamRef:       input.Ref,
	})
	if err != nil {
		return nil, mapTeamError(err, "list team members failed")
	}

	return &listMembersOutput{Body: listMembersOutputBody{Members: newTeamMemberResponses(members)}}, nil
}

type listMembersOutput struct {
	Body listMembersOutputBody
}

type listMembersOutputBody struct {
	Members []teamMemberResponse `json:"members"`
}

func newTeamMemberResponses(members []domainteam.Member) []teamMemberResponse {
	responses := make([]teamMemberResponse, 0, len(members))
	for _, member := range members {
		responses = append(responses, newTeamMemberResponse(member))
	}

	return responses
}
