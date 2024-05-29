package auth
import (
	"anonymous/commons"
	"anonymous/helpers"
	"anonymous/models"
	"anonymous/types"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	MustInsert(tx *sqlx.Tx, user *models.User) error
	CheckDuplicates(email string) (string, error)
	CheckDuplicatesU(username string) (string, error)
	GetUserDataByID(id string) (*models.LoggedInUser, error)
	ChangePassword(password, id string) error
	ToggleStatus(users []string, status bool) error
	GetAllUsersData() (*[]models.LoggedInUser, error)
	SetContactVerified(userId string) error
	GetUser(string, string) (*models.User, error)
}

type AuthService struct {
	users      UserRepo
	txProvider types.TxProvider
	logger     types.Logger
	jwt        types.JWTProvider
}

func Service(
	users UserRepo,
	txProvider types.TxProvider,
	logger types.Logger,
	jwt types.JWTProvider,
) *AuthService {
	return &AuthService{
		users:      users,
		txProvider: txProvider,
		logger:     logger,
		jwt:        jwt,
	}
}

func (s AuthService) Register(data *registrationPayload) (*string, *models.LoggedInUser, error) {
	hash, err := helpers.Hash(data.Password)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	res, err := s.users.CheckDuplicatesU(data.Username)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, nil, commons.Errors.InternalServerError
		}
		if res == "username" {
			return nil, nil, commons.Errors.DuplicateUsername
		}
	
	result, err := s.users.CheckDuplicates(data.Email)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	if result == "email" {
		return nil, nil, commons.Errors.DuplicateEmail
	}

	user := &models.User{
		ID:             uuid.NewString(),
		Username:       data.Username,
		Password:       hash,
		Email:          data.Email,
		EmailVerified:  false,
		JoinedAt:       time.Now().UTC(),
		Active:         true,
		ProfilePicture: "",
	}
	tx, err := s.txProvider.Provide()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	err = s.users.MustInsert(tx, user)
	if err != nil {
		s.logger.Error(err.Error())
		err = tx.Rollback()
		if err != nil {
			s.logger.Error(fmt.Sprintf(
				"Error while rolling back transaction: %s",
				err,
			))
		}
		return nil, nil, commons.Errors.InternalServerError
	}
	err = tx.Commit()
	if err != nil {
		s.logger.Error(fmt.Sprintf(
			"Error while committing transaction: %s",
			err,
		))
		return nil, nil, commons.Errors.InternalServerError
	}
	token, err := s.jwt.Encode(map[string]interface{}{
		"id": user.ID,
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf(
			"Error while encoding token: %s",
			err,
		))
		return nil, nil, commons.Errors.TokenEncodingFailed
	}
	userData, err := s.users.GetUserDataByID(user.ID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	return &token, userData, nil
}


func (s *AuthService) Login(data *loginPayload) (*string, *models.LoggedInUser, error) {
	var err error
	var user *models.User
	var lookupErr string
	switch data.Method {
	case "username":
		user, err = s.users.GetUser(data.Method, data.Username)
		lookupErr = commons.Codes.UsernameNotFound
	case "email":
		user, err = s.users.GetUser(data.Method, data.Email)
		lookupErr = commons.Codes.EmailNotFound
	}
	if err != nil {
		if errors.Is(err, commons.Errors.ResourceNotFound) {
			return nil, nil, types.ServiceError{
				StatusCode: http.StatusBadRequest,
				ErrorCode:  lookupErr,
			}
		}
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	if !helpers.HashMatchesString(user.Password, data.Password) {
		return nil, nil, types.ServiceError{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  commons.Codes.WrongPassword,
		}
	}
	token, err := s.jwt.Encode(map[string]interface{}{
		"id": user.ID,
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf(
			"Error while encoding token: %s",
			err,
		))
		return nil, nil, commons.Errors.TokenEncodingFailed
	}
	userData, err := s.users.GetUserDataByID(user.ID)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	return &token, userData, nil
}

