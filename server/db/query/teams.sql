-- name: CreateTeam :one
INSERT INTO teams (name, slug, created_by_user_id)
VALUES ($1, $2, $3)
RETURNING id, name, slug, created_by_user_id, created_at, updated_at, deleted_at;

-- name: CreateTeamMember :one
WITH inserted AS (
    INSERT INTO team_members (team_id, user_id, role)
    VALUES ($1, $2, $3)
    RETURNING id, team_id, user_id, role, created_at, updated_at
)
SELECT inserted.id,
       inserted.team_id,
       inserted.user_id,
       users.email,
       inserted.role,
       inserted.created_at,
       inserted.updated_at
FROM inserted
JOIN users ON users.id = inserted.user_id;

-- name: ListTeamsForUser :many
SELECT teams.id, teams.name, teams.slug, teams.created_by_user_id, teams.created_at, teams.updated_at, teams.deleted_at
FROM teams
JOIN team_members
    ON team_members.team_id = teams.id
    AND team_members.user_id = $1
    AND team_members.deleted_at IS NULL
WHERE teams.deleted_at IS NULL
ORDER BY teams.created_at DESC, teams.id DESC;

-- name: GetTeamForUser :one
SELECT teams.id, teams.name, teams.slug, teams.created_by_user_id, teams.created_at, teams.updated_at, teams.deleted_at
FROM teams
JOIN team_members
    ON team_members.team_id = teams.id
    AND team_members.user_id = $2
    AND team_members.deleted_at IS NULL
WHERE teams.id = $1
  AND teams.deleted_at IS NULL;

-- name: GetActiveTeamMemberRole :one
SELECT team_members.role
FROM team_members
JOIN teams ON teams.id = team_members.team_id
WHERE team_members.team_id = $1
  AND team_members.user_id = $2
  AND team_members.deleted_at IS NULL
  AND teams.deleted_at IS NULL;

-- name: UpdateTeam :one
UPDATE teams
SET name = $2,
    slug = $3
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, name, slug, created_by_user_id, created_at, updated_at, deleted_at;

-- name: SoftDeleteTeam :one
UPDATE teams
SET deleted_at = now()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id;

-- name: ListActiveTeamMembers :many
SELECT team_members.id,
       team_members.team_id,
       team_members.user_id,
       users.email,
       team_members.role,
       team_members.created_at,
       team_members.updated_at
FROM team_members
JOIN users ON users.id = team_members.user_id
JOIN teams ON teams.id = team_members.team_id
WHERE team_members.team_id = $1
  AND team_members.deleted_at IS NULL
  AND teams.deleted_at IS NULL
ORDER BY team_members.created_at ASC, team_members.id ASC;

-- name: GetActiveTeamMember :one
SELECT team_members.id,
       team_members.team_id,
       team_members.user_id,
       users.email,
       team_members.role,
       team_members.created_at,
       team_members.updated_at
FROM team_members
JOIN users ON users.id = team_members.user_id
JOIN teams ON teams.id = team_members.team_id
WHERE team_members.team_id = $1
  AND team_members.user_id = $2
  AND team_members.deleted_at IS NULL
  AND teams.deleted_at IS NULL;

-- name: UpdateTeamMemberRole :one
WITH updated AS (
    UPDATE team_members
    SET role = $3
    WHERE team_id = $1
      AND user_id = $2
      AND deleted_at IS NULL
    RETURNING id, team_id, user_id, role, created_at, updated_at
)
SELECT updated.id,
       updated.team_id,
       updated.user_id,
       users.email,
       updated.role,
       updated.created_at,
       updated.updated_at
FROM updated
JOIN users ON users.id = updated.user_id;

-- name: CountActiveTeamOwners :one
SELECT count(*)::int4
FROM team_members
WHERE team_id = $1
  AND role = 'owner'
  AND deleted_at IS NULL;
