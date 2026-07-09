package repositories

import "errors"

// ErrNotFound is returned by any repository lookup that finds no
// matching row.
var ErrNotFound = errors.New("record not found")
