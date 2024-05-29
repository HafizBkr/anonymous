package users

import (
	"anonymous/models"
	"anonymous/types"
	"anonymous/utils"
	"encoding/json"
	"net/http"
)

type IUserService interface {
	ChangePassword(data *changePasswordPayload, userData *models.LoggedInUser) error
	ToggleUserAccountStatus(users []string, status bool) error
	GetAllUsersData() (*[]models.LoggedInUser, error)
}

type UsersHandler struct {
	service IUserService
	logger  types.Logger
}

func Handler(
	service IUserService,
	logger types.Logger,
) *UsersHandler {
	return &UsersHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UsersHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	currUser, ok := r.Context().Value("user").(*models.LoggedInUser)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	payload := &changePasswordPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		utils.HandleBodyDecodingErr(w, err, h.logger)
		return
	}
	errs := payload.Validate()
	if errs != nil {
		utils.WriteValidationError(w, errs)
    return
	}
	err = h.service.ChangePassword(payload, currUser)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteData(w, http.StatusOK, nil)
}

func (h *UsersHandler) HandleToggleStatus(w http.ResponseWriter, r *http.Request) {
	payload := &toggleUserStatusPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		utils.HandleBodyDecodingErr(w, err, h.logger)
		return
	}
	errs := payload.Validate()
	if errs != nil {
		utils.WriteValidationError(w, errs)
    return
	}
	err = h.service.ToggleUserAccountStatus(payload.IDs, payload.Active)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteData(w, http.StatusOK, nil)
}


func (h *UsersHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAllUsersData()
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	utils.WriteData(w, http.StatusOK, map[string]interface{}{
		"users": *data,
	})
}
