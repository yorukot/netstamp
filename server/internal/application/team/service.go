package team

import (
	"context"
	"regexp"
	"strings"

	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)

func IsValidSlug(value string) bool {
	return slugPattern.MatchString(value)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTeam(ctx context.Context, input CreateTeamInput) (domainteam.Team, error) {
	name, err := normalizeRequired(input.Name)
	if err != nil {
		return domainteam.Team{}, err
	}
	slug, err := normalizeSlug(input.Slug)
	if err != nil {
		return domainteam.Team{}, err
	}

	return s.repo.CreateTeamWithOwner(ctx, CreateTeamStorageInput{
		Name:            name,
		Slug:            slug,
		CreatedByUserID: input.CurrentUserID,
	})
}

func (s *Service) ListTeams(ctx context.Context, input ListTeamsInput) ([]domainteam.Team, error) {
	return s.repo.ListTeamsForUser(ctx, input.CurrentUserID)
}

func (s *Service) GetTeam(ctx context.Context, input GetTeamInput) (domainteam.Team, error) {
	return s.repo.GetTeamForUser(ctx, input.TeamRef, input.CurrentUserID)
}

func (s *Service) UpdateTeam(ctx context.Context, input UpdateTeamInput) (domainteam.Team, error) {
	team, err := s.repo.GetTeamForUser(ctx, input.TeamRef, input.CurrentUserID)
	if err != nil {
		return domainteam.Team{}, err
	}

	role, err := s.repo.GetMemberRole(ctx, team.ID, input.CurrentUserID)
	if err != nil {
		return domainteam.Team{}, err
	}
	if !isManager(role) {
		return domainteam.Team{}, ErrForbidden
	}

	name := team.Name
	slug := team.Slug
	if input.Name != nil {
		name, err = normalizeRequired(*input.Name)
		if err != nil {
			return domainteam.Team{}, err
		}
	}
	if input.Slug != nil {
		slug, err = normalizeSlug(*input.Slug)
		if err != nil {
			return domainteam.Team{}, err
		}
	}
	if input.Name == nil && input.Slug == nil {
		return domainteam.Team{}, ErrInvalidInput
	}

	return s.repo.UpdateTeam(ctx, UpdateTeamStorageInput{
		TeamID: team.ID,
		Name:   name,
		Slug:   slug,
	})
}

func (s *Service) DeleteTeam(ctx context.Context, input DeleteTeamInput) error {
	team, err := s.repo.GetTeamForUser(ctx, input.TeamRef, input.CurrentUserID)
	if err != nil {
		return err
	}

	role, err := s.repo.GetMemberRole(ctx, team.ID, input.CurrentUserID)
	if err != nil {
		return err
	}
	if role != domainteam.RoleOwner {
		return ErrForbidden
	}

	return s.repo.SoftDeleteTeam(ctx, team.ID)
}

func (s *Service) ListMembers(ctx context.Context, input ListMembersInput) ([]domainteam.Member, error) {
	team, err := s.repo.GetTeamForUser(ctx, input.TeamRef, input.CurrentUserID)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetMemberRole(ctx, team.ID, input.CurrentUserID); err != nil {
		return nil, err
	}

	return s.repo.ListMembers(ctx, team.ID)
}

func (s *Service) AddMember(ctx context.Context, input AddMemberInput) (domainteam.Member, error) {
	team, err := s.repo.GetTeamForUser(ctx, input.TeamRef, input.CurrentUserID)
	if err != nil {
		return domainteam.Member{}, err
	}

	actorRole, err := s.repo.GetMemberRole(ctx, team.ID, input.CurrentUserID)
	if err != nil {
		return domainteam.Member{}, err
	}
	if !isManager(actorRole) {
		return domainteam.Member{}, ErrForbidden
	}
	if err := validateRole(input.Role); err != nil {
		return domainteam.Member{}, err
	}
	if !canAssignRole(actorRole, input.Role) {
		return domainteam.Member{}, ErrForbidden
	}

	return s.repo.AddMember(ctx, AddMemberStorageInput{
		TeamID: team.ID,
		UserID: input.UserID,
		Role:   input.Role,
	})
}

func (s *Service) UpdateMemberRole(ctx context.Context, input UpdateMemberRoleInput) (domainteam.Member, error) {
	team, err := s.repo.GetTeamForUser(ctx, input.TeamRef, input.CurrentUserID)
	if err != nil {
		return domainteam.Member{}, err
	}

	actorRole, err := s.repo.GetMemberRole(ctx, team.ID, input.CurrentUserID)
	if err != nil {
		return domainteam.Member{}, err
	}
	if !isManager(actorRole) {
		return domainteam.Member{}, ErrForbidden
	}
	if err := validateRole(input.Role); err != nil {
		return domainteam.Member{}, err
	}
	if !canAssignRole(actorRole, input.Role) {
		return domainteam.Member{}, ErrForbidden
	}

	member, err := s.repo.GetMember(ctx, team.ID, input.UserID)
	if err != nil {
		return domainteam.Member{}, err
	}
	if actorRole == domainteam.RoleAdmin && member.Role == domainteam.RoleOwner {
		return domainteam.Member{}, ErrForbidden
	}
	if member.Role == domainteam.RoleOwner && input.Role != domainteam.RoleOwner {
		owners, err := s.repo.CountOwners(ctx, team.ID)
		if err != nil {
			return domainteam.Member{}, err
		}
		if owners <= 1 {
			return domainteam.Member{}, ErrLastOwner
		}
	}

	return s.repo.UpdateMemberRole(ctx, UpdateMemberRoleStorageInput{
		TeamID: team.ID,
		UserID: input.UserID,
		Role:   input.Role,
	})
}

func normalizeRequired(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidInput
	}

	return value, nil
}

func normalizeSlug(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || !IsValidSlug(value) {
		return "", ErrInvalidInput
	}

	return value, nil
}

func validateRole(role domainteam.Role) error {
	switch role {
	case domainteam.RoleOwner, domainteam.RoleAdmin, domainteam.RoleEditor, domainteam.RoleViewer:
		return nil
	default:
		return ErrInvalidRole
	}
}

func isManager(role domainteam.Role) bool {
	return role == domainteam.RoleOwner || role == domainteam.RoleAdmin
}

func canAssignRole(actorRole domainteam.Role, targetRole domainteam.Role) bool {
	switch actorRole {
	case domainteam.RoleOwner:
		return targetRole == domainteam.RoleAdmin || targetRole == domainteam.RoleEditor || targetRole == domainteam.RoleViewer
	case domainteam.RoleAdmin:
		return targetRole == domainteam.RoleEditor || targetRole == domainteam.RoleViewer
	default:
		return false
	}
}
