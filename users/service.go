package users

import (
	"anonymous/auth"
	"anonymous/commons"
	"anonymous/helpers"
	"anonymous/models"
	"anonymous/types"
)

type UserService struct {
	users      auth.UserRepo
	txProvider types.TxProvider
	logger     types.Logger
}

func Service(
	users auth.UserRepo,
	txProvider types.TxProvider,
	logger types.Logger,
) *UserService {
	return &UserService{
		users:      users,
		txProvider: txProvider,
		logger:     logger,
	}
}

func (s *UserService) ChangePassword(
	data *changePasswordPayload,
	userData *models.LoggedInUser,
) error {
	if !helpers.HashMatchesString(userData.Password, data.Old) {
		return commons.Errors.InternalServerError
	}
	hash, err := helpers.Hash(data.New)
	if err != nil {
		s.logger.Error(err.Error())
		return commons.Errors.InternalServerError
	}
	err = s.users.ChangePassword(hash, userData.ID)
	if err != nil {
		s.logger.Error(err.Error())
		return commons.Errors.InternalServerError
	}
	return nil
}

func (s *UserService) ToggleUserAccountStatus(users []string, status bool) error {
	err := s.users.ToggleStatus(users, status)
	if err != nil {
		s.logger.Error(err.Error())
		return commons.Errors.InternalServerError
	}
	return nil
}

func (s *UserService) GetAllUsersData() (*[]models.LoggedInUser, error) {
	data, err := s.users.GetAllUsersData()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, commons.Errors.InternalServerError
	}
	return data, nil
}
