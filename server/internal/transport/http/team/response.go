package team

import (
	"time"

	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

type teamOutput struct {
	Body teamOutputBody
}

type teamOutputBody struct {
	Team teamResponse `json:"team"`
}

type memberOutput struct {
	Body memberOutputBody
}

type memberOutputBody struct {
	Member teamMemberResponse `json:"member"`
}

type teamResponse struct {
	ID              string    `json:"id" format:"uuid"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	CreatedByUserID string    `json:"createdByUserId" format:"uuid"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type teamMemberResponse struct {
	ID        string    `json:"id" format:"uuid"`
	TeamID    string    `json:"teamId" format:"uuid"`
	UserID    string    `json:"userId" format:"uuid"`
	Email     string    `json:"email" format:"email"`
	Role      string    `json:"role" enum:"owner,admin,editor,viewer"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newTeamResponse(team domainteam.Team) teamResponse {
	return teamResponse{
		ID:              team.ID,
		Name:            team.Name,
		Slug:            team.Slug,
		CreatedByUserID: team.CreatedByUserID,
		CreatedAt:       team.CreatedAt,
		UpdatedAt:       team.UpdatedAt,
	}
}

func newTeamMemberResponse(member domainteam.Member) teamMemberResponse {
	return teamMemberResponse{
		ID:        member.ID,
		TeamID:    member.TeamID,
		UserID:    member.UserID,
		Email:     member.Email,
		Role:      string(member.Role),
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}
