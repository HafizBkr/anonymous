package replies

import (
    "anonymous/auth"
    "anonymous/commons"
    "anonymous/models"
    "anonymous/validator"
)

type CommentReplyService interface {
    CreateCommentReply(token string, payload *CommentReplyPayload) (*models.CommentReply, error)
    GetCommentRepliesByCommentID(commentID string) ([]*models.CommentReply, error)
    GetCommentReply(id string) (*models.CommentReply, error)
    UpdateCommentReply(token string, id string, payload *UpdateCommentReplyPayload) (*models.CommentReply, error)
    DeleteCommentReply(token string, id string) error
}

type commentReplyService struct {
    repo        CommentReplyRepo
    authService auth.AuthService
}

func NewCommentReplyService(repo CommentReplyRepo, authService auth.AuthService) CommentReplyService {
    return &commentReplyService{
        repo:        repo,
        authService: authService,
    }
}

func (s *commentReplyService) CreateCommentReply(token string, payload *CommentReplyPayload) (*models.CommentReply, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }
    payload.UserID = userID

    if err := validator.ValidateStruct(payload); err != nil {
        return nil, err
    }

    return s.repo.CreateCommentReply(payload)
}

func (s *commentReplyService) GetCommentRepliesByCommentID(commentID string) ([]*models.CommentReply, error) {
    return s.repo.GetCommentRepliesByCommentID(commentID)
}

func (s *commentReplyService) GetCommentReply(id string) (*models.CommentReply, error) {
    return s.repo.GetCommentReply(id)
}

func (s *commentReplyService) UpdateCommentReply(token string, id string, payload *UpdateCommentReplyPayload) (*models.CommentReply, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, commons.Errors.AuthenticationFailed
    }

    if err := validator.ValidateStruct(payload); err != nil {
        return nil, err
    }

    return s.repo.UpdateCommentReply(id, payload, userID)
}

func (s *commentReplyService) DeleteCommentReply(token string, id string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return commons.Errors.AuthenticationFailed
    }

    return s.repo.DeleteCommentReply(userID, id)
}
