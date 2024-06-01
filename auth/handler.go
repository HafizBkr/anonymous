package auth

import (
	"anonymous/commons"
	"anonymous/models"
	"anonymous/utils"
	"encoding/json"
	"net/http"
	"anonymous/types"
)

type IAuthService interface {
	Register(*registrationPayload) (*string, *models.LoggedInUser, error)
	Login(*loginPayload) (*string, *models.LoggedInUser, error)
	  VerifyEmail(token string) error  
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

func (h *AuthHandler) HandleEmailVerification(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    if token == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Token manquant"))
        return
    }

    err := h.service.VerifyEmail(token)
    if err != nil {
        if serr, ok := err.(types.ServiceError); ok {
            // Si l'erreur est de type ServiceError, nous pouvons extraire le statut HTTP et le code d'erreur
            utils.WriteServiceError(w, serr)
            return
        }
        // Si ce n'est pas une erreur de service spécifique, nous pouvons simplement renvoyer une erreur interne du serveur
        utils.WriteServiceError(w, types.ServiceError{
            StatusCode: http.StatusInternalServerError,
            ErrorCode:  "InternalError",
        })
        return
    }

    // Rediriger vers la page verified.html après la vérification réussie
    http.Redirect(w, r, "/static/verified.html", http.StatusSeeOther)
}



