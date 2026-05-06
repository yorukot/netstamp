package team

import "errors"

var (
	ErrTeamNotFound          = errors.New("team not found")
	ErrTeamSlugAlreadyExists = errors.New("team slug already exists")
	ErrForbidden             = errors.New("team action forbidden")
	ErrInvalidInput          = errors.New("team input invalid")
	ErrInvalidRole           = errors.New("team member role invalid")
	ErrMemberAlreadyExists   = errors.New("team member already exists")
	ErrMemberNotFound        = errors.New("team member not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrLastOwner             = errors.New("team must keep an owner")
)
