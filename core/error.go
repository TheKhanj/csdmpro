package core

import "errors"

var (
	ERR_PLAYER_NOT_FOUND   error = errors.New("player not found")
	ERR_NO_CANDIDATE_FOUND error = errors.New("no candidate found")
)
