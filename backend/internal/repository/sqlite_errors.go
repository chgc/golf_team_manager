package repository

import (
	"errors"

	"modernc.org/sqlite"
)

func isSQLiteConstraintError(err error) bool {
	var sqliteError *sqlite.Error
	if errors.As(err, &sqliteError) {
		return int(sqliteError.Code())&0xFF == 19
	}

	return false
}
