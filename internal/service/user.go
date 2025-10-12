package service

import (
	"PattyWagon/internal/constants"
	"PattyWagon/internal/model"
	"PattyWagon/internal/utils"
	"context"
	"strings"
)

func (s *Service) UsernameLogin(ctx context.Context, username string, password string, role int16) (token string, err error) {
	user, err := s.repository.SelectUserCredentialsByUsernameAndRole(ctx, username, role)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", constants.ErrUserNotFound
		}
		return "", err
	}
	if user.ID == 0 { // user not found
		return "", constants.ErrUserNotFound
	}

	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", constants.ErrUserWrongPassword
	}

	token, err = utils.GenerateToken(user.ID, role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Register(ctx context.Context, userReq model.User, password string, role int16) (string, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user, err := s.repository.InsertUser(ctx, userReq, passwordHash)
	if err != nil {
		if utils.IsErrDBConstraint(err) {
			return "", constants.ErrDuplicate
		}
		return "", err
	}

	token, err := utils.GenerateToken(user.ID, role)
	if err != nil {
		return "", err
	}
	return token, nil
}
