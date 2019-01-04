package postgres

import "errors"

var ErrNoRows = errors.New("no records returned")
var ErrMultipleRows = errors.New("multiple records returned")
