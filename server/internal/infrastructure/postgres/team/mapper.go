package pgteam

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	domainteam "github.com/yorukot/netstamp/internal/domain/team"
	"github.com/yorukot/netstamp/internal/infrastructure/postgres/sqlc"
)

func mapTeam(row sqlc.Team) domainteam.Team {
	return domainteam.Team{
		ID:              row.ID.String(),
		Name:            row.Name,
		Slug:            row.Slug,
		CreatedByUserID: row.CreatedByUserID.String(),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		DeletedAt:       timePtr(row.DeletedAt),
	}
}

func mapCreateMember(row sqlc.CreateTeamMemberRow) domainteam.Member {
	return domainteam.Member{
		ID:        row.ID.String(),
		TeamID:    row.TeamID.String(),
		UserID:    row.UserID.String(),
		Email:     row.Email,
		Role:      domainteam.Role(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func timePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

func mapListMember(row sqlc.ListActiveTeamMembersRow) domainteam.Member {
	return domainteam.Member{
		ID:        row.ID.String(),
		TeamID:    row.TeamID.String(),
		UserID:    row.UserID.String(),
		Email:     row.Email,
		Role:      domainteam.Role(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func mapGetMember(row sqlc.GetActiveTeamMemberRow) domainteam.Member {
	return domainteam.Member{
		ID:        row.ID.String(),
		TeamID:    row.TeamID.String(),
		UserID:    row.UserID.String(),
		Email:     row.Email,
		Role:      domainteam.Role(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func mapUpdateMember(row sqlc.UpdateTeamMemberRoleRow) domainteam.Member {
	return domainteam.Member{
		ID:        row.ID.String(),
		TeamID:    row.TeamID.String(),
		UserID:    row.UserID.String(),
		Email:     row.Email,
		Role:      domainteam.Role(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
