package repository

import (
	"WeenieHut/internal/model"
	"context"
	"database/sql"
	"errors"
	"strconv"
)

const (
	insertUserQuery                   = `INSERT INTO users (email, phone, password_hash) VALUES ($1, $2, $3) RETURNING id`
	selectUserCredentialsByEmailQuery = `SELECT id, phone, password_hash FROM users WHERE email = $1`
	selectUserCredentialsByPhoneQuery = `SELECT id, email, password_hash FROM users WHERE phone = $1`
	selectUserById                    = `SELECT * FROM users WHERE id = $1`
	checkUserExistById                = `SELECT 1 FROM users WHERE id = $1 LIMIT 1`
	checkEmailExistQuery              = `SELECT 1 FROM users WHERE email = $1 AND id != $2 LIMIT 1`
	checkPhoneExistQuery              = `SELECT 1 FROM users WHERE phone = $1 AND id != $2 LIMIT 1`
	updateUserQuery                   = `
		UPDATE users 
		SET 
			file_id = COALESCE($1, file_id),
			file_uri = COALESCE($2, file_uri),
			file_thumbnail_uri = COALESCE($3, file_thumbnail_uri),
			bank_account_name = COALESCE($4, bank_account_name),
			bank_account_holder = COALESCE($5, bank_account_holder),
			bank_account_number = COALESCE($6, bank_account_number),
			email = COALESCE($7, email),
			phone = COALESCE($8, phone),
			updated_at = NOW()
		WHERE id = $9
		RETURNING *
	`
)

func (q *Queries) InsertUser(ctx context.Context, user model.User, passwordHash string) (res model.User, err error) {
	err = q.db.QueryRowContext(ctx, insertUserQuery,
		user.Email,
		user.Phone,
		passwordHash,
	).Scan(&user.ID)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (q *Queries) SelectUserCredentialsByEmail(ctx context.Context, email string) (res model.User, err error) {
	err = q.db.QueryRowContext(ctx, selectUserCredentialsByEmailQuery, email).Scan(&res.ID, &res.Phone, &res.PasswordHash)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

func (q *Queries) SelectUserCredentialsByPhone(ctx context.Context, email string) (res model.User, err error) {
	err = q.db.QueryRowContext(ctx, selectUserCredentialsByPhoneQuery, email).Scan(&res.ID, &res.Email, &res.PasswordHash)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

type UpdateUserProfileParams struct {
	FileID            int64
	FileThumbnailURI  string
	FileURI           string
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber string
	Email             string
	Phone             string
	UserID            int64
}

func (q *Queries) UpdateUserProfile(ctx context.Context, params UpdateUserProfileParams) (model.User, error) {
	query := updateUserQuery

	// Convert empty strings to nil for proper COALESCE behavior
	var fileID, fileURI, fileThumbnailURI, bankAccountName, bankAccountHolder, bankAccountNumber, email, phone interface{}

	if params.FileID == 0 {
		fileID = nil
	} else {
		fileID = strconv.Itoa(int(params.FileID))
	}

	if params.FileURI == "" {
		fileURI = nil
	} else {
		fileURI = params.FileURI
	}

	if params.FileThumbnailURI == "" {
		fileThumbnailURI = nil
	} else {
		fileThumbnailURI = params.FileThumbnailURI
	}

	if params.BankAccountName == "" {
		bankAccountName = nil
	} else {
		bankAccountName = params.BankAccountName
	}

	if params.BankAccountHolder == "" {
		bankAccountHolder = nil
	} else {
		bankAccountHolder = params.BankAccountHolder
	}

	if params.BankAccountNumber == "" {
		bankAccountNumber = nil
	} else {
		bankAccountNumber = params.BankAccountNumber
	}

	if params.Email == "" {
		email = nil
	} else {
		email = params.Email
	}

	if params.Phone == "" {
		phone = nil
	} else {
		phone = params.Phone
	}

	row := q.db.QueryRowContext(ctx, query,
		fileID,
		fileURI,
		fileThumbnailURI,
		bankAccountName,
		bankAccountHolder,
		bankAccountNumber,
		email,
		phone,
		params.UserID,
	)

	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.PasswordHash,
		&u.Name,
		&u.FileID,
		&u.FileURI,
		&u.FileThumbnailURI,
		&u.BankAccountName,
		&u.BankAccountHolder,
		&u.BankAccountNumber,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return model.User{}, err
	}

	return u, nil
}

func (q *Queries) GetUserProfile(ctx context.Context, id int64) (model.User, error) {
	query := selectUserById
	row := q.db.QueryRowContext(ctx, query, id)
	var u model.User
	if err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.PasswordHash,
		&u.Name,
		&u.FileID,
		&u.FileURI,
		&u.FileThumbnailURI,
		&u.BankAccountName,
		&u.BankAccountHolder,
		&u.BankAccountNumber,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return u, nil
}

func (q *Queries) IsUserExist(ctx context.Context, userID int64) (bool, error) {
	query := checkUserExistById

	var exists int
	err := q.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (q *Queries) IsEmailExist(ctx context.Context, email string, excludeUserID int64) (bool, error) {
	query := checkEmailExistQuery

	var exists int
	err := q.db.QueryRowContext(ctx, query, email, excludeUserID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (q *Queries) IsPhoneExist(ctx context.Context, phone string, excludeUserID int64) (bool, error) {
	query := checkPhoneExistQuery

	var exists int
	err := q.db.QueryRowContext(ctx, query, phone, excludeUserID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
