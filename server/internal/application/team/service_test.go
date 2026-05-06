package team

import (
	"context"
	"errors"
	"testing"
	"time"

	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

func TestCreateTeamNormalizesInputAndCreatesOwnerMembership(t *testing.T) {
	repo := &fakeTeamRepository{
		createdTeam: domainteam.Team{ID: "team-1"},
	}
	service := NewService(repo)

	_, err := service.CreateTeam(context.Background(), CreateTeamInput{
		CurrentUserID: "user-1",
		Name:          "  Engineering  ",
		Slug:          "  platform-team  ",
	})
	if err != nil {
		t.Fatalf("create team: %v", err)
	}

	if repo.gotCreateInput.Name != "Engineering" {
		t.Fatalf("expected trimmed name, got %q", repo.gotCreateInput.Name)
	}
	if repo.gotCreateInput.Slug != "platform-team" {
		t.Fatalf("expected trimmed slug, got %q", repo.gotCreateInput.Slug)
	}
	if repo.gotCreateInput.CreatedByUserID != "user-1" {
		t.Fatalf("expected owner user id, got %q", repo.gotCreateInput.CreatedByUserID)
	}
}

func TestCreateTeamRejectsInvalidSlug(t *testing.T) {
	repo := &fakeTeamRepository{}
	service := NewService(repo)

	_, err := service.CreateTeam(context.Background(), CreateTeamInput{
		CurrentUserID: "user-1",
		Name:          "Engineering",
		Slug:          "Platform_Team",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
	if repo.gotCreateInput.Slug != "" {
		t.Fatalf("expected create not to be called, got %#v", repo.gotCreateInput)
	}
}

func TestDeleteTeamRequiresOwner(t *testing.T) {
	repo := &fakeTeamRepository{actorRole: domainteam.RoleAdmin}
	service := NewService(repo)

	err := service.DeleteTeam(context.Background(), DeleteTeamInput{
		CurrentUserID: "admin-user",
		TeamID:        "team-1",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	if repo.gotSoftDeleteTeamID != "" {
		t.Fatalf("expected delete not to be called, got %q", repo.gotSoftDeleteTeamID)
	}
}

func TestAddMemberRoleRestrictions(t *testing.T) {
	tests := []struct {
		name      string
		actorRole domainteam.Role
		newRole   domainteam.Role
		wantErr   error
	}{
		{
			name:      "owner cannot add owner",
			actorRole: domainteam.RoleOwner,
			newRole:   domainteam.RoleOwner,
			wantErr:   ErrForbidden,
		},
		{
			name:      "admin cannot add admin",
			actorRole: domainteam.RoleAdmin,
			newRole:   domainteam.RoleAdmin,
			wantErr:   ErrForbidden,
		},
		{
			name:      "owner can add admin",
			actorRole: domainteam.RoleOwner,
			newRole:   domainteam.RoleAdmin,
		},
		{
			name:      "admin can add viewer",
			actorRole: domainteam.RoleAdmin,
			newRole:   domainteam.RoleViewer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeTeamRepository{
				actorRole:   tt.actorRole,
				addedMember: domainteam.Member{ID: "member-1", Role: tt.newRole},
			}
			service := NewService(repo)

			_, err := service.AddMember(context.Background(), AddMemberInput{
				CurrentUserID: "actor-user",
				TeamID:        "team-1",
				UserID:        "target-user",
				Role:          tt.newRole,
			})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil && repo.gotAddMember.UserID != "" {
				t.Fatalf("expected add not to be called, got %#v", repo.gotAddMember)
			}
			if tt.wantErr == nil && repo.gotAddMember.Role != tt.newRole {
				t.Fatalf("expected add role %q, got %q", tt.newRole, repo.gotAddMember.Role)
			}
		})
	}
}

func TestUpdateMemberRoleRestrictions(t *testing.T) {
	tests := []struct {
		name       string
		actorRole  domainteam.Role
		memberRole domainteam.Role
		newRole    domainteam.Role
		owners     int
		wantErr    error
	}{
		{
			name:       "owner cannot change anyone to owner",
			actorRole:  domainteam.RoleOwner,
			memberRole: domainteam.RoleViewer,
			newRole:    domainteam.RoleOwner,
			owners:     1,
			wantErr:    ErrForbidden,
		},
		{
			name:       "admin cannot change anyone to admin",
			actorRole:  domainteam.RoleAdmin,
			memberRole: domainteam.RoleViewer,
			newRole:    domainteam.RoleAdmin,
			owners:     1,
			wantErr:    ErrForbidden,
		},
		{
			name:       "admin cannot change owner",
			actorRole:  domainteam.RoleAdmin,
			memberRole: domainteam.RoleOwner,
			newRole:    domainteam.RoleViewer,
			owners:     2,
			wantErr:    ErrForbidden,
		},
		{
			name:       "cannot remove last owner",
			actorRole:  domainteam.RoleOwner,
			memberRole: domainteam.RoleOwner,
			newRole:    domainteam.RoleAdmin,
			owners:     1,
			wantErr:    ErrLastOwner,
		},
		{
			name:       "owner can change member to admin",
			actorRole:  domainteam.RoleOwner,
			memberRole: domainteam.RoleViewer,
			newRole:    domainteam.RoleAdmin,
			owners:     1,
		},
		{
			name:       "admin can change member to editor",
			actorRole:  domainteam.RoleAdmin,
			memberRole: domainteam.RoleViewer,
			newRole:    domainteam.RoleEditor,
			owners:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeTeamRepository{
				actorRole: tt.actorRole,
				member: domainteam.Member{
					ID:     "member-1",
					UserID: "target-user",
					Role:   tt.memberRole,
				},
				owners:        tt.owners,
				updatedMember: domainteam.Member{ID: "member-1", Role: tt.newRole},
			}
			service := NewService(repo)

			_, err := service.UpdateMemberRole(context.Background(), UpdateMemberRoleInput{
				CurrentUserID: "actor-user",
				TeamID:        "team-1",
				UserID:        "target-user",
				Role:          tt.newRole,
			})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil && repo.gotUpdateMemberRole.UserID != "" {
				t.Fatalf("expected update not to be called, got %#v", repo.gotUpdateMemberRole)
			}
			if tt.wantErr == nil && repo.gotUpdateMemberRole.Role != tt.newRole {
				t.Fatalf("expected update role %q, got %q", tt.newRole, repo.gotUpdateMemberRole.Role)
			}
		})
	}
}

type fakeTeamRepository struct {
	createdTeam         domainteam.Team
	gotCreateInput      CreateTeamStorageInput
	teams               []domainteam.Team
	team                domainteam.Team
	actorRole           domainteam.Role
	roleErr             error
	updatedTeam         domainteam.Team
	gotUpdateTeam       UpdateTeamStorageInput
	gotSoftDeleteTeamID string
	members             []domainteam.Member
	member              domainteam.Member
	addedMember         domainteam.Member
	gotAddMember        AddMemberStorageInput
	updatedMember       domainteam.Member
	gotUpdateMemberRole UpdateMemberRoleStorageInput
	owners              int
}

func (r *fakeTeamRepository) CreateTeamWithOwner(_ context.Context, input CreateTeamStorageInput) (domainteam.Team, error) {
	r.gotCreateInput = input
	if r.createdTeam.ID != "" {
		return r.createdTeam, nil
	}
	return domainteam.Team{
		ID:              "team-1",
		Name:            input.Name,
		Slug:            input.Slug,
		CreatedByUserID: input.CreatedByUserID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func (r *fakeTeamRepository) ListTeamsForUser(context.Context, string) ([]domainteam.Team, error) {
	return r.teams, nil
}

func (r *fakeTeamRepository) GetTeamForUser(context.Context, string, string) (domainteam.Team, error) {
	if r.team.ID != "" {
		return r.team, nil
	}
	return domainteam.Team{ID: "team-1", Name: "Engineering", Slug: "engineering"}, nil
}

func (r *fakeTeamRepository) GetMemberRole(context.Context, string, string) (domainteam.Role, error) {
	if r.roleErr != nil {
		return "", r.roleErr
	}
	if r.actorRole != "" {
		return r.actorRole, nil
	}
	return domainteam.RoleOwner, nil
}

func (r *fakeTeamRepository) UpdateTeam(_ context.Context, input UpdateTeamStorageInput) (domainteam.Team, error) {
	r.gotUpdateTeam = input
	if r.updatedTeam.ID != "" {
		return r.updatedTeam, nil
	}
	return domainteam.Team{ID: input.TeamID, Name: input.Name, Slug: input.Slug}, nil
}

func (r *fakeTeamRepository) SoftDeleteTeam(_ context.Context, teamID string) error {
	r.gotSoftDeleteTeamID = teamID
	return nil
}

func (r *fakeTeamRepository) ListMembers(context.Context, string) ([]domainteam.Member, error) {
	return r.members, nil
}

func (r *fakeTeamRepository) GetMember(context.Context, string, string) (domainteam.Member, error) {
	if r.member.ID != "" {
		return r.member, nil
	}
	return domainteam.Member{ID: "member-1", UserID: "target-user", Role: domainteam.RoleViewer}, nil
}

func (r *fakeTeamRepository) AddMember(_ context.Context, input AddMemberStorageInput) (domainteam.Member, error) {
	r.gotAddMember = input
	if r.addedMember.ID != "" {
		return r.addedMember, nil
	}
	return domainteam.Member{ID: "member-1", TeamID: input.TeamID, UserID: input.UserID, Role: input.Role}, nil
}

func (r *fakeTeamRepository) UpdateMemberRole(_ context.Context, input UpdateMemberRoleStorageInput) (domainteam.Member, error) {
	r.gotUpdateMemberRole = input
	if r.updatedMember.ID != "" {
		return r.updatedMember, nil
	}
	return domainteam.Member{ID: "member-1", TeamID: input.TeamID, UserID: input.UserID, Role: input.Role}, nil
}

func (r *fakeTeamRepository) CountOwners(context.Context, string) (int, error) {
	if r.owners > 0 {
		return r.owners, nil
	}
	return 1, nil
}
