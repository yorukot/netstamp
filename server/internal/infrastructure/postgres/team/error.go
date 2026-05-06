package pgteam

import (
	"fmt"

	appteam "github.com/yorukot/netstamp/internal/application/team"
	"github.com/yorukot/netstamp/internal/infrastructure/postgres"
)

func mapCreateTeamError(err error) error {
	if postgres.IsUniqueViolation(err, "uq_teams_slug") {
		return fmt.Errorf("team slug already exists: %w", appteam.ErrTeamSlugAlreadyExists)
	}
	if postgres.IsForeignKeyViolation(err, "teams_created_by_user_id_fkey") {
		return fmt.Errorf("user not found: %w", appteam.ErrUserNotFound)
	}

	return err
}

func mapCreateTeamMemberError(err error) error {
	if postgres.IsUniqueViolation(err, "uq_team_members_active_team_user") {
		return fmt.Errorf("team member already exists: %w", appteam.ErrMemberAlreadyExists)
	}
	if postgres.IsForeignKeyViolation(err, "team_members_team_id_fkey") {
		return fmt.Errorf("team not found: %w", appteam.ErrTeamNotFound)
	}
	if postgres.IsForeignKeyViolation(err, "team_members_user_id_fkey") {
		return fmt.Errorf("user not found: %w", appteam.ErrUserNotFound)
	}

	return err
}
