package utils

import (
	"database/sql"
	"strconv"
)

func ToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func ToInt(ni sql.NullString) int64 {
	if ni.Valid {
		if i, err := strconv.ParseInt(ni.String, 10, 64); err == nil {
			return i
		}
	}
	return 0
}
