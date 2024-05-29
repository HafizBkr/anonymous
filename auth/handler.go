package auth

import (
	"anonymous/commons"
	"anonymous/models"
	"anonymous/utils"
	"encoding/json"
	"net/http"
)

type IAuthService interface {
	Register(*registrationPayload) (*string, *models.LoggedInUser, error)
	Login(*loginPayload) (*string, *models.LoggedInUser, error)
}

type AuthHandler struct {
	service IAuthService
	logger  commons.Logger
}

func NewAuthHandler(
	service IAuthService,
	logger commons.Logger,
) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuthHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	payload := &registrationPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		utils.HandleBodyDecodingErr(w, err, h.logger)
		return
	}
	validationErr := payload.Validate()
	if validationErr != nil {
		utils.WriteValidationError(w, validationErr)
		return
	}
	token, userData, err := h.service.Register(payload)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	data := map[string]interface{}{
		"token": *token,
		"user":  *userData,
	}
	utils.WriteData(w, http.StatusCreated, data)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	payload := &loginPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		utils.HandleBodyDecodingErr(w, err, h.logger)
		return
	}
	validationErr := payload.Validate()
	if validationErr != nil {
		utils.WriteValidationError(w, validationErr)
		return
	}
	token, userData, err := h.service.Login(payload)
	if err != nil {
		utils.WriteServiceError(w, err)
		return
	}
	data := map[string]interface{}{
		"token": *token,
		"user":  *userData,
	}
	utils.WriteData(w, http.StatusOK, data)
}

func (h *AuthHandler) HandleGetCurrentUserData(w http.ResponseWriter, r *http.Request) {
	currUser, ok := r.Context().Value("user").(*models.LoggedInUser)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	utils.WriteData(w, http.StatusOK, map[string]interface{}{
		"user": *currUser,
	})
	return
}
