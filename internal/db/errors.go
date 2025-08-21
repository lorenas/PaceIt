package db

import "errors"

var ErrMissingDSN = errors.New("DATABASE_URL not set")