package comments

import (
    "anonymous/auth"
    "anonymous/commons"
    "anonymous/models"
    "anonymous/validator"
)

type CommentService interface {
    CreateComment(token string, payload *CommentPayload) (*models.Comment, error)
    GetComment(id string) (*models.Comment, error)
    UpdateComment(token string,  id string, payload *UpdateCommentPayload) (*models.Comment, error)
    DeleteComment(token string, id string) error
    GetCommentsByPostID(postID string) ([]*models.Comment, error)
    GetCommentsCountByPostID(postID string) (int, error)
    AddOrUpdateReaction(token, commentID, reactionType string) (*models.CommentReaction, error)
    CountReactions(commentID string) (map[string]int, error)
}


type commentService struct {
    repo        CommentRepo
    authService auth.AuthService
}

func NewCommentService(repo CommentRepo, authService auth.AuthService) CommentService {
    return &commentService{
        repo:        repo,
        authService: authService,
    }
}
func (s *commentService) GetCommentsByPostID(postID string) ([]*models.Comment, error) {
    return s.repo.GetCommentsByPostID(postID)
}


func (s *commentService) CreateComment(token string, payload *CommentPayload) (*models.Comment, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }
    payload.UserID = userID

    if err := validator.ValidateStruct(payload); err != nil {
        return nil, err
    }

    return s.repo.CreateComment(payload)
}

func (s *commentService) GetComment(id string) (*models.Comment, error) {
    return s.repo.GetComment(id)
}

func (s *commentService) UpdateComment(token string, id string, payload *UpdateCommentPayload) (*models.Comment, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }

    if err := validator.ValidateStruct(payload); err != nil {
        return nil, err
    }

    return s.repo.UpdateComment(id, payload, userID)
}

func (s *commentService) DeleteComment(token string, commentID string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return commons.Errors.AuthenticationFailed
    }

    return s.repo.DeleteComment(userID, commentID)
}

func (s *commentService) GetCommentsCountByPostID(postID string) (int, error) {
    return s.repo.GetCommentsCountByPostID(postID)
}

func (s *commentService) AddOrUpdateReaction(token, commentID, reactionType string) (*models.CommentReaction, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }

    return s.repo.AddOrUpdateReaction(commentID, userID, reactionType)
}


func (s *commentService) CountReactions(commentID string) (map[string]int, error) {
    return s.repo.CountReactions(commentID)
}