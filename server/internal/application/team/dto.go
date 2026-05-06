package team

import domainteam "github.com/yorukot/netstamp/internal/domain/team"

type CreateTeamInput struct {
	CurrentUserID string
	Name          string
	Slug          string
}

type CreateTeamStorageInput struct {
	Name            string
	Slug            string
	CreatedByUserID string
}

type ListTeamsInput struct {
	CurrentUserID string
}

type GetTeamInput struct {
	CurrentUserID string
	TeamID        string
}

type UpdateTeamInput struct {
	CurrentUserID string
	TeamID        string
	Name          *string
	Slug          *string
}

type UpdateTeamStorageInput struct {
	TeamID string
	Name   string
	Slug   string
}

type DeleteTeamInput struct {
	CurrentUserID string
	TeamID        string
}

type ListMembersInput struct {
	CurrentUserID string
	TeamID        string
}

type AddMemberInput struct {
	CurrentUserID string
	TeamID        string
	UserID        string
	Role          domainteam.Role
}

type AddMemberStorageInput struct {
	TeamID string
	UserID string
	Role   domainteam.Role
}

type UpdateMemberRoleInput struct {
	CurrentUserID string
	TeamID        string
	UserID        string
	Role          domainteam.Role
}

type UpdateMemberRoleStorageInput struct {
	TeamID string
	UserID string
	Role   domainteam.Role
}
