package utils

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/lib/pq"
)

// Database error checking functions
func IsErrDBConstraint(err error) bool {
	return strings.Contains(err.Error(), "unique constraint")
}

func IsDBError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common database errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return true
	}

	// Check for standard database errors
	if errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, sql.ErrTxDone) ||
		errors.Is(err, sql.ErrConnDone) {
		return true
	}

	// Check for driver errors
	if errors.Is(err, driver.ErrBadConn) {
		return true
	}

	// Check error message patterns
	errMsg := strings.ToLower(err.Error())
	dbErrorPatterns := []string{
		"connection refused",
		"connection reset",
		"connection timeout",
		"database is locked",
		"constraint violation",
		"foreign key constraint",
		"unique constraint",
		"check constraint",
		"syntax error",
		"relation does not exist",
		"column does not exist",
		"duplicate key",
		"deadlock detected",
	}

	for _, pattern := range dbErrorPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

func GetDBErrorType(err error) string {
	if err == nil {
		return ""
	}

	// PostgreSQL specific errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			return "unique_constraint"
		case "23503": // foreign_key_violation
			return "foreign_key_constraint"
		case "23514": // check_violation
			return "check_constraint"
		case "42P01": // undefined_table
			return "table_not_found"
		case "42703": // undefined_column
			return "column_not_found"
		case "08006": // connection_failure
			return "connection_error"
		case "40001": // serialization_failure
			return "deadlock"
		default:
			return "database_error"
		}
	}

	// Standard SQL errors
	if errors.Is(err, sql.ErrNoRows) {
		return "not_found"
	}
	if errors.Is(err, sql.ErrTxDone) {
		return "transaction_done"
	}
	if errors.Is(err, sql.ErrConnDone) {
		return "connection_done"
	}

	return "unknown_database_error"
}

func IsAppError(err error) bool {
	return !IsDBError(err)
}
