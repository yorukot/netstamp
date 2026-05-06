package team

import (
	"context"

	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

type Repository interface {
	CreateTeamWithOwner(ctx context.Context, input CreateTeamStorageInput) (domainteam.Team, error)
	ListTeamsForUser(ctx context.Context, userID string) ([]domainteam.Team, error)
	GetTeamForUser(ctx context.Context, teamRef string, userID string) (domainteam.Team, error)
	GetMemberRole(ctx context.Context, teamID string, userID string) (domainteam.Role, error)
	UpdateTeam(ctx context.Context, input UpdateTeamStorageInput) (domainteam.Team, error)
	SoftDeleteTeam(ctx context.Context, teamID string) error
	ListMembers(ctx context.Context, teamID string) ([]domainteam.Member, error)
	GetMember(ctx context.Context, teamID string, userID string) (domainteam.Member, error)
	AddMember(ctx context.Context, input AddMemberStorageInput) (domainteam.Member, error)
	UpdateMemberRole(ctx context.Context, input UpdateMemberRoleStorageInput) (domainteam.Member, error)
	CountOwners(ctx context.Context, teamID string) (int, error)
}
