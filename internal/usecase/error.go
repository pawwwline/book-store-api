package usecase

import "errors"

var (
	ErrDbInfrastructure error = errors.New("database infrastructure error")
	ErrCache            error = errors.New("cache error")
)
