package posts

import (
    "anonymous/auth"
    "anonymous/commons"
    "anonymous/models"
    "anonymous/validator"
)

type PostService interface {
    CreatePost(token string, payload *PostPayload) (*models.Post, error)
    GetPost(id string) (*models.Post, error)
     UpdatePost(token string, id string, payload *PostPayload) (*models.Post, error) // Ajout du paramètre "id"
    DeletePost(token string, id string) error
}

type postService struct {
    repo        PostRepo
    authService auth.AuthService
}

func NewPostService(repo PostRepo, authService auth.AuthService) PostService {
    return &postService{
        repo:        repo,
        authService: authService,
    }
}

func (s *postService) CreatePost(token string, payload *PostPayload) (*models.Post, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }
    payload.UserID = userID

    if err := validator.ValidateStruct(payload); err != nil {
        return nil, err
    }

    return s.repo.CreatePost(payload)
}

func (s *postService) GetPost(id string) (*models.Post, error) {
    return s.repo.GetPost(id)
}

func (s *postService) UpdatePost(token string, id string, payload *PostPayload) (*models.Post, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }
    payload.UserID = userID

    if err := validator.ValidateStruct(payload); err != nil {
        return nil, err
    }

    return s.repo.UpdatePost(id, payload) // Correction de l'appel de fonction pour inclure "id"
}


func (s *postService) DeletePost(token string, id string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return commons.Errors.AuthenticationFailed
    }

    return s.repo.DeletePost(userID, id)
}
