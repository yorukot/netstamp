package pgteam

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	appteam "github.com/yorukot/netstamp/internal/application/team"
	domainteam "github.com/yorukot/netstamp/internal/domain/team"
	"github.com/yorukot/netstamp/internal/infrastructure/postgres"
	"github.com/yorukot/netstamp/internal/infrastructure/postgres/sqlc"
)

type TeamRepository struct {
	queries *sqlc.Queries
	tx      *postgres.Transactor
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{
		queries: sqlc.New(pool),
		tx:      postgres.NewTransactor(pool),
	}
}

func (r *TeamRepository) CreateTeamWithOwner(ctx context.Context, input appteam.CreateTeamStorageInput) (domainteam.Team, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "teams", "postgres.teams.create_with_owner", "INSERT", "INSERT teams and owner membership")
	defer span.End()

	userID, err := postgres.ParseUUID(input.CreatedByUserID, appteam.ErrUserNotFound)
	if err != nil {
		return domainteam.Team{}, err
	}

	var team domainteam.Team
	err = r.tx.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		q := r.queries.WithTx(tx)

		row, err := q.CreateTeam(ctx, sqlc.CreateTeamParams{
			Name:            input.Name,
			Slug:            input.Slug,
			CreatedByUserID: userID,
		})
		if err != nil {
			return mapCreateTeamError(err)
		}
		team = mapTeam(row)

		if _, err := q.CreateTeamMember(ctx, sqlc.CreateTeamMemberParams{
			TeamID: row.ID,
			UserID: userID,
			Role:   sqlc.TeamMemberRoleOwner,
		}); err != nil {
			return mapCreateTeamMemberError(err)
		}

		return nil
	})
	if err != nil {
		postgres.RecordDBSpanError(span, err)
		return domainteam.Team{}, err
	}

	return team, nil
}

func (r *TeamRepository) ListTeamsForUser(ctx context.Context, userIDValue string) ([]domainteam.Team, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "teams", "postgres.teams.list_for_user", "SELECT", "SELECT teams for member")
	defer span.End()

	userID, err := postgres.ParseUUID(userIDValue, appteam.ErrTeamNotFound)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.ListTeamsForUser(ctx, userID)
	if err != nil {
		postgres.RecordDBSpanError(span, err)
		return nil, err
	}

	teams := make([]domainteam.Team, 0, len(rows))
	for _, row := range rows {
		teams = append(teams, mapTeam(row))
	}

	return teams, nil
}

func (r *TeamRepository) GetTeamForUser(ctx context.Context, teamIDValue string, userIDValue string) (domainteam.Team, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "teams", "postgres.teams.select_for_user", "SELECT", "SELECT team for member")
	defer span.End()

	teamID, userID, err := parseTeamAndUserIDs(teamIDValue, userIDValue)
	if err != nil {
		return domainteam.Team{}, err
	}

	row, err := r.queries.GetTeamForUser(ctx, sqlc.GetTeamForUserParams{
		ID:     teamID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainteam.Team{}, appteam.ErrTeamNotFound
		}
		postgres.RecordDBSpanError(span, err)
		return domainteam.Team{}, err
	}

	return mapTeam(row), nil
}

func (r *TeamRepository) GetMemberRole(ctx context.Context, teamIDValue string, userIDValue string) (domainteam.Role, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "team_members", "postgres.team_members.select_role", "SELECT", "SELECT active team member role")
	defer span.End()

	teamID, userID, err := parseTeamAndUserIDs(teamIDValue, userIDValue)
	if err != nil {
		return "", err
	}

	role, err := r.queries.GetActiveTeamMemberRole(ctx, sqlc.GetActiveTeamMemberRoleParams{
		TeamID: teamID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", appteam.ErrTeamNotFound
		}
		postgres.RecordDBSpanError(span, err)
		return "", err
	}

	return domainteam.Role(role), nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, input appteam.UpdateTeamStorageInput) (domainteam.Team, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "teams", "postgres.teams.update", "UPDATE", "UPDATE team")
	defer span.End()

	teamID, err := postgres.ParseUUID(input.TeamID, appteam.ErrTeamNotFound)
	if err != nil {
		return domainteam.Team{}, err
	}

	row, err := r.queries.UpdateTeam(ctx, sqlc.UpdateTeamParams{
		ID:   teamID,
		Name: input.Name,
		Slug: input.Slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainteam.Team{}, appteam.ErrTeamNotFound
		}
		if mapped := mapCreateTeamError(err); mapped != err {
			return domainteam.Team{}, mapped
		}
		postgres.RecordDBSpanError(span, err)
		return domainteam.Team{}, err
	}

	return mapTeam(row), nil
}

func (r *TeamRepository) SoftDeleteTeam(ctx context.Context, teamIDValue string) error {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "teams", "postgres.teams.soft_delete", "UPDATE", "SOFT DELETE team")
	defer span.End()

	teamID, err := postgres.ParseUUID(teamIDValue, appteam.ErrTeamNotFound)
	if err != nil {
		return err
	}

	if _, err := r.queries.SoftDeleteTeam(ctx, teamID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appteam.ErrTeamNotFound
		}
		postgres.RecordDBSpanError(span, err)
		return err
	}

	return nil
}

func (r *TeamRepository) ListMembers(ctx context.Context, teamIDValue string) ([]domainteam.Member, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "team_members", "postgres.team_members.list", "SELECT", "SELECT active team members")
	defer span.End()

	teamID, err := postgres.ParseUUID(teamIDValue, appteam.ErrTeamNotFound)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.ListActiveTeamMembers(ctx, teamID)
	if err != nil {
		postgres.RecordDBSpanError(span, err)
		return nil, err
	}

	members := make([]domainteam.Member, 0, len(rows))
	for _, row := range rows {
		members = append(members, mapListMember(row))
	}

	return members, nil
}

func (r *TeamRepository) GetMember(ctx context.Context, teamIDValue string, userIDValue string) (domainteam.Member, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "team_members", "postgres.team_members.select", "SELECT", "SELECT active team member")
	defer span.End()

	teamID, userID, err := parseTeamAndUserIDs(teamIDValue, userIDValue)
	if err != nil {
		return domainteam.Member{}, appteam.ErrMemberNotFound
	}

	row, err := r.queries.GetActiveTeamMember(ctx, sqlc.GetActiveTeamMemberParams{
		TeamID: teamID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainteam.Member{}, appteam.ErrMemberNotFound
		}
		postgres.RecordDBSpanError(span, err)
		return domainteam.Member{}, err
	}

	return mapGetMember(row), nil
}

func (r *TeamRepository) AddMember(ctx context.Context, input appteam.AddMemberStorageInput) (domainteam.Member, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "team_members", "postgres.team_members.insert", "INSERT", "INSERT team member")
	defer span.End()

	teamID, userID, err := parseTeamAndUserIDs(input.TeamID, input.UserID)
	if err != nil {
		return domainteam.Member{}, err
	}

	row, err := r.queries.CreateTeamMember(ctx, sqlc.CreateTeamMemberParams{
		TeamID: teamID,
		UserID: userID,
		Role:   sqlc.TeamMemberRole(input.Role),
	})
	if err != nil {
		mapped := mapCreateTeamMemberError(err)
		if mapped == err {
			postgres.RecordDBSpanError(span, err)
		}
		return domainteam.Member{}, mapped
	}

	return mapCreateMember(row), nil
}

func (r *TeamRepository) UpdateMemberRole(ctx context.Context, input appteam.UpdateMemberRoleStorageInput) (domainteam.Member, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "team_members", "postgres.team_members.update_role", "UPDATE", "UPDATE team member role")
	defer span.End()

	teamID, userID, err := parseTeamAndUserIDs(input.TeamID, input.UserID)
	if err != nil {
		return domainteam.Member{}, appteam.ErrMemberNotFound
	}

	row, err := r.queries.UpdateTeamMemberRole(ctx, sqlc.UpdateTeamMemberRoleParams{
		TeamID: teamID,
		UserID: userID,
		Role:   sqlc.TeamMemberRole(input.Role),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainteam.Member{}, appteam.ErrMemberNotFound
		}
		postgres.RecordDBSpanError(span, err)
		return domainteam.Member{}, err
	}

	return mapUpdateMember(row), nil
}

func (r *TeamRepository) CountOwners(ctx context.Context, teamIDValue string) (int, error) {
	ctx, span := postgres.StartDBSpan(ctx, pgteamTracer, "team_members", "postgres.team_members.count_owners", "SELECT", "COUNT active team owners")
	defer span.End()

	teamID, err := postgres.ParseUUID(teamIDValue, appteam.ErrTeamNotFound)
	if err != nil {
		return 0, err
	}

	count, err := r.queries.CountActiveTeamOwners(ctx, teamID)
	if err != nil {
		postgres.RecordDBSpanError(span, err)
		return 0, err
	}

	return int(count), nil
}

func parseTeamAndUserIDs(teamIDValue string, userIDValue string) (uuid.UUID, uuid.UUID, error) {
	teamID, err := postgres.ParseUUID(teamIDValue, appteam.ErrTeamNotFound)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	userID, err := postgres.ParseUUID(userIDValue, appteam.ErrUserNotFound)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return teamID, userID, nil
}
