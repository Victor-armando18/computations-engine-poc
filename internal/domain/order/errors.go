package order

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrRuleViolation = errors.New("rule violation")
	ErrInvalidPatch  = errors.New("invalid patch")
	ErrIdempotentHit = errors.New("idempotent replay")
)
