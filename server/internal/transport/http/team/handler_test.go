package team

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"

	appauth "github.com/yorukot/netstamp/internal/application/auth"
	appteam "github.com/yorukot/netstamp/internal/application/team"
	domainteam "github.com/yorukot/netstamp/internal/domain/team"
)

const (
	testUserID   = "11111111-1111-1111-1111-111111111111"
	testTeamID   = "22222222-2222-2222-2222-222222222222"
	testMemberID = "33333333-3333-3333-3333-333333333333"
)

func TestCreateTeamReturnsCreatedTeam(t *testing.T) {
	_, api := humatest.New(t)
	repo := &handlerTeamRepository{}
	NewHandler(appteam.NewService(repo), &handlerTokenVerifier{
		claims: appauth.AccessTokenClaims{Subject: testUserID, Email: "user@example.com"},
	}).RegisterRoutes(api)

	res := api.Post("/teams", map[string]any{
		"name": " Engineering ",
		"slug": "engineering",
	}, "Authorization: Bearer valid-token")

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", res.Code)
	}

	var body teamOutputBody
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Team.ID != testTeamID {
		t.Fatalf("expected team id, got %q", body.Team.ID)
	}
	if body.Team.Slug != "engineering" {
		t.Fatalf("expected slug, got %q", body.Team.Slug)
	}
	if repo.gotCreateInput.CreatedByUserID != testUserID {
		t.Fatalf("expected current user id, got %q", repo.gotCreateInput.CreatedByUserID)
	}
}

func TestCreateTeamRejectsInvalidSlugPattern(t *testing.T) {
	_, api := humatest.New(t)
	NewHandler(appteam.NewService(&handlerTeamRepository{}), &handlerTokenVerifier{
		claims: appauth.AccessTokenClaims{Subject: testUserID, Email: "user@example.com"},
	}).RegisterRoutes(api)

	res := api.Post("/teams", map[string]any{
		"name": "Engineering",
		"slug": "Engineering_Team",
	}, "Authorization: Bearer valid-token")

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", res.Code)
	}
}

func TestCreateTeamRequiresBearerToken(t *testing.T) {
	_, api := humatest.New(t)
	NewHandler(appteam.NewService(&handlerTeamRepository{}), &handlerTokenVerifier{}).RegisterRoutes(api)

	res := api.Post("/teams", map[string]any{
		"name": "Engineering",
		"slug": "engineering",
	})

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", res.Code)
	}
}

func TestGetTeamAcceptsSlugRef(t *testing.T) {
	_, api := humatest.New(t)
	repo := &handlerTeamRepository{}
	NewHandler(appteam.NewService(repo), &handlerTokenVerifier{
		claims: appauth.AccessTokenClaims{Subject: testUserID, Email: "user@example.com"},
	}).RegisterRoutes(api)

	res := api.Get("/teams/engineering", "Authorization: Bearer valid-token")

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	if repo.gotTeamRef != "engineering" {
		t.Fatalf("expected slug ref, got %q", repo.gotTeamRef)
	}
}

func TestDeleteTeamMapsNonOwnerToForbidden(t *testing.T) {
	_, api := humatest.New(t)
	NewHandler(appteam.NewService(&handlerTeamRepository{
		actorRole: domainteam.RoleAdmin,
	}), &handlerTokenVerifier{
		claims: appauth.AccessTokenClaims{Subject: testUserID, Email: "user@example.com"},
	}).RegisterRoutes(api)

	res := api.Delete("/teams/engineering", "Authorization: Bearer valid-token")

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", res.Code)
	}
}

func TestAddMemberRejectsOwnerRoleAsForbidden(t *testing.T) {
	_, api := humatest.New(t)
	NewHandler(appteam.NewService(&handlerTeamRepository{
		actorRole: domainteam.RoleOwner,
	}), &handlerTokenVerifier{
		claims: appauth.AccessTokenClaims{Subject: testUserID, Email: "user@example.com"},
	}).RegisterRoutes(api)

	res := api.Post("/teams/engineering/members", map[string]any{
		"userId": testMemberID,
		"role":   "owner",
	}, "Authorization: Bearer valid-token")

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", res.Code)
	}
}

type handlerTokenVerifier struct {
	claims appauth.AccessTokenClaims
	err    error
}

func (v *handlerTokenVerifier) VerifyAccessToken(context.Context, string) (appauth.AccessTokenClaims, error) {
	if v.err != nil {
		return appauth.AccessTokenClaims{}, v.err
	}
	return v.claims, nil
}

type handlerTeamRepository struct {
	gotCreateInput      appteam.CreateTeamStorageInput
	gotTeamRef          string
	actorRole           domainteam.Role
	gotSoftDeleteTeamID string
}

func (r *handlerTeamRepository) CreateTeamWithOwner(_ context.Context, input appteam.CreateTeamStorageInput) (domainteam.Team, error) {
	r.gotCreateInput = input
	return domainteam.Team{
		ID:              testTeamID,
		Name:            input.Name,
		Slug:            input.Slug,
		CreatedByUserID: input.CreatedByUserID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func (r *handlerTeamRepository) ListTeamsForUser(context.Context, string) ([]domainteam.Team, error) {
	return []domainteam.Team{{
		ID:              testTeamID,
		Name:            "Engineering",
		Slug:            "engineering",
		CreatedByUserID: testUserID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}}, nil
}

func (r *handlerTeamRepository) GetTeamForUser(_ context.Context, teamRef string, _ string) (domainteam.Team, error) {
	r.gotTeamRef = teamRef
	return domainteam.Team{
		ID:              testTeamID,
		Name:            "Engineering",
		Slug:            "engineering",
		CreatedByUserID: testUserID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func (r *handlerTeamRepository) GetMemberRole(context.Context, string, string) (domainteam.Role, error) {
	if r.actorRole != "" {
		return r.actorRole, nil
	}
	return domainteam.RoleOwner, nil
}

func (r *handlerTeamRepository) UpdateTeam(_ context.Context, input appteam.UpdateTeamStorageInput) (domainteam.Team, error) {
	return domainteam.Team{
		ID:              input.TeamID,
		Name:            input.Name,
		Slug:            input.Slug,
		CreatedByUserID: testUserID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func (r *handlerTeamRepository) SoftDeleteTeam(_ context.Context, teamID string) error {
	r.gotSoftDeleteTeamID = teamID
	return nil
}

func (r *handlerTeamRepository) ListMembers(context.Context, string) ([]domainteam.Member, error) {
	return []domainteam.Member{{ID: testMemberID, TeamID: testTeamID, UserID: testUserID, Email: "user@example.com", Role: domainteam.RoleOwner}}, nil
}

func (r *handlerTeamRepository) GetMember(context.Context, string, string) (domainteam.Member, error) {
	return domainteam.Member{ID: testMemberID, TeamID: testTeamID, UserID: testMemberID, Email: "member@example.com", Role: domainteam.RoleViewer}, nil
}

func (r *handlerTeamRepository) AddMember(_ context.Context, input appteam.AddMemberStorageInput) (domainteam.Member, error) {
	return domainteam.Member{ID: testMemberID, TeamID: input.TeamID, UserID: input.UserID, Email: "member@example.com", Role: input.Role}, nil
}

func (r *handlerTeamRepository) UpdateMemberRole(_ context.Context, input appteam.UpdateMemberRoleStorageInput) (domainteam.Member, error) {
	return domainteam.Member{ID: testMemberID, TeamID: input.TeamID, UserID: input.UserID, Email: "member@example.com", Role: input.Role}, nil
}

func (r *handlerTeamRepository) CountOwners(context.Context, string) (int, error) {
	return 1, nil
}
