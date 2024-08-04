package points

import (
    "anonymous/commons"
    "anonymous/types"
    "fmt"
)

// PointsService définit l'interface pour le service de points.
type PointsService interface {
    LikeUserProfile(userID, likedBy string) error
    GetUserProfileLikes(userID string) (int, error)
    DecodeToken(token string) (map[string]interface{}, error)
}

// pointsService implémente l'interface PointsService.
type pointsService struct {
    repo        PointsRepo
    logger      types.Logger
    jwtProvider types.JWTProvider
}

// NewPointsService crée une nouvelle instance de pointsService avec les dépendances fournies.
func NewPointsService(repo PointsRepo, logger types.Logger, jwtProvider types.JWTProvider) PointsService {
    return &pointsService{
        repo:        repo,
        logger:      logger,
        jwtProvider: jwtProvider,
    }
}

// LikeUserProfile permet à un utilisateur de liker le profil d'un autre utilisateur.
func (s *pointsService) LikeUserProfile(userID, likedBy string) error {
    err := s.repo.LikeUserProfile(userID, likedBy)
    if err != nil {
        s.logger.Error(err.Error())
        return commons.Errors.DuplicateLike
    }
    return nil
}

// GetUserProfileLikes récupère le nombre de likes pour un profil utilisateur donné.
func (s *pointsService) GetUserProfileLikes(userID string) (int, error) {
    if userID == "" {
        s.logger.Error("UserID is empty")
        return 0, commons.Errors.BadRequest
    }

    count, err := s.repo.CountUserProfileLikes(userID)
    if err != nil {
        s.logger.Error(err.Error())
        return 0, commons.Errors.InternalServerError
    }
    return count, nil
}

// DecodeToken décode le token JWT et retourne les claims.
func (s *pointsService) DecodeToken(token string) (map[string]interface{}, error) {
    claims, err := s.jwtProvider.Decode(token)
    if err != nil {
        return nil, fmt.Errorf("Error decoding token: %w", err)
    }
    return claims, nil
}
