package auth

import (
	"anonymous/commons"
	"anonymous/emails"
	"anonymous/helpers"
	"anonymous/models"
	"anonymous/types"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
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
	VerifyEmail(token string) error 
	SetEmailVerificationToken(userID, token string) error
	FindByVerificationToken(token string) (*models.User, error)
	Update(user *models.User) error 
	SetPasswordResetToken(email, token string) error
	FindByPasswordResetToken(token string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdatePassword(userID, password string) error
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
	resUsername, err := s.users.CheckDuplicatesU(data.Username)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	if resUsername == "username" {
		return nil, nil, commons.Errors.DuplicateUsername
	}

	resEmail, err := s.users.CheckDuplicates(data.Email)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	if resEmail == "email" {
		return nil, nil, commons.Errors.DuplicateEmail
	}

	hash, err := helpers.Hash(data.Password)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	verificationToken := uuid.New().String()

   
	user := &models.User{
		ID:             uuid.NewString(),
		Username:       data.Username,
		Password:       hash,
		Email:          data.Email, 
		EmailVerified:  false,
		JoinedAt:       time.Now().UTC(),
		Active:         true,
		ProfilePicture: "",
		EmailVerificationToken: verificationToken,
	}
	
 // Enregistrement du token de vérification dans la base de données
    if err := s.users.SetEmailVerificationToken(user.ID, verificationToken); err != nil {
        s.logger.Error(err.Error())
        return nil, nil, commons.Errors.InternalServerError
    }
    
    email := os.Getenv("EMAIL")
    password := os.Getenv("PASSWORD")
       sender := emails.NewSender(email, password)
        sender.SendVerificationEmail([]string{data.Email}, user.EmailVerificationToken)


	tx, err := s.txProvider.Provide()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				s.logger.Error(fmt.Sprintf(
					"Error while rolling back transaction: %s",
					rollbackErr,
				))
			}
		}
	}()

	if err := s.users.MustInsert(tx, user); err != nil {
		s.logger.Error(err.Error())
		return nil, nil, commons.Errors.InternalServerError
	}
	
	if err := tx.Commit(); err != nil {
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
        user, err = s.users.GetUser("username", data.Username)
        lookupErr = commons.Codes.UsernameNotFound
    case "email":
        // Utilisation de GetUserByEmail qui gère déjà le cryptage
        user, err = s.users.GetUserByEmail(data.Email)
        lookupErr = commons.Codes.EmailNotFound
    default:
        return nil, nil, commons.Errors.InvalidLoginMethod
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

    // Vérifier si l'e-mail de l'utilisateur est vérifié
    if !user.EmailVerified {
        return nil, nil, types.ServiceError{
            StatusCode: http.StatusUnauthorized,
            ErrorCode:  commons.Codes.EmailNotVerified,
        }
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

    // Récupérer les données de l'utilisateur avec l'email décrypté
    userData, err := s.users.GetUserDataByID(user.ID)
    if err != nil {
        s.logger.Error(err.Error())
        return nil, nil, commons.Errors.InternalServerError
    }

    return &token, userData, nil
}

func (s *AuthService) VerifyEmail(token string) error {
	err := s.users.VerifyEmail(token)
	if err != nil {
		if errors.Is(err, commons.Errors.ResourceNotFound) {
			return types.ServiceError{
				StatusCode: http.StatusBadRequest,
				ErrorCode:  commons.Codes.InvalidVerificationToken,
			}
		}
		s.logger.Error(err.Error())
		return commons.Errors.InternalServerError
	}
	return nil
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
    // Strip the "Bearer " prefix if it exists
    if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
        tokenString = tokenString[7:]
    }
    
    userID, err := s.jwt.ValidateToken(tokenString)
    if err != nil {
        return "", err
    }
    return userID, nil
}

// ForgotPassword initiates the password reset process for a user
func (s *AuthService) ForgotPassword(email string) error {
	 // Generate a reset token
	 resetToken := uuid.New().String()
    
	 // Store the reset token in the database
	 if err := s.users.SetPasswordResetToken(email, resetToken); err != nil {
		 s.logger.Error(err.Error())
		 return commons.Errors.InternalServerError
	 }
    // Find user by email to ensure they exist
    _, err := s.users.GetUserByEmail(email)
    if err != nil {
        if errors.Is(err, commons.Errors.ResourceNotFound) {
            return nil
        }
        s.logger.Error(err.Error())
        return commons.Errors.InternalServerError
    }
    
    // Send the password reset email
    emailAddr := os.Getenv("EMAIL")
    password := os.Getenv("PASSWORD")
    sender := emails.NewSender(emailAddr, password)
    
    if err := sender.SendPasswordResetEmail([]string{email}, resetToken); err != nil {
        s.logger.Error(fmt.Sprintf("Error sending password reset email: %s", err))
        return commons.Errors.InternalServerError
    }
    
    return nil
}

// ResetPassword completes the password reset process
func (s *AuthService) ResetPassword(data *resetPasswordPayload) error {
    // Find user by reset token
    user, err := s.users.FindByPasswordResetToken(data.Token)
    if err != nil {
        if errors.Is(err, commons.Errors.ResourceNotFound) {
            return types.ServiceError{
                StatusCode: http.StatusBadRequest,
                ErrorCode:  commons.Codes.InvalidResetToken,
            }
        }
        s.logger.Error(err.Error())
        return commons.Errors.InternalServerError
    }
    
    // Hash the new password
    hash, err := helpers.Hash(data.NewPassword)
    if err != nil {
        s.logger.Error(err.Error())
        return commons.Errors.InternalServerError
    }
    
    // Update the user's password
    if err := s.users.UpdatePassword(user.ID, hash); err != nil {
        s.logger.Error(err.Error())
        return commons.Errors.InternalServerError
    }
    
    return nil
}