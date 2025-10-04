package repository

import (
	"PattyWagon/internal/model"
	"context"
)

const (
	insertUserQuery                      = `INSERT INTO users (email, username, role, password_hash) VALUES ($1, $2, $3, $4) RETURNING id`
	selectUserCredentialsByUsernameQuery = `SELECT id, email, password_hash FROM users WHERE username = $1 AND role = $2`
)

func (q *Queries) InsertUser(ctx context.Context, user model.User, passwordHash string) (res model.User, err error) {
	err = q.db.QueryRowContext(ctx, insertUserQuery,
		user.Email,
		user.Username,
		user.Role,
		passwordHash,
	).Scan(&user.ID)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (q *Queries) SelectUserCredentialsByUsernameAndRole(ctx context.Context, username string, role int16) (res model.User, err error) {
	err = q.db.QueryRowContext(ctx, selectUserCredentialsByUsernameQuery, username, role).Scan(&res.ID, &res.Email, &res.PasswordHash)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}
