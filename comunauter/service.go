package comunauter

import (
	"anonymous/auth"
	"anonymous/models"
	"fmt"
	 "github.com/google/uuid"
)

type CommunityService interface {
	CreateCommunity(payload *CommunityPayload, token string) (*models.Community, error)
	GetCommunity(id string) (*models.Community, error)
	GetAllCommunities() ([]*models.Community, error)
	JoinCommunity(token string, communityID string) error
	GetCommunityUsers(communityID string) ([]*models.User, error)
}

type communityService struct {
	repo        CommunityRepo
	authService auth.AuthService
}

func NewCommunityService(repo CommunityRepo, authService auth.AuthService) CommunityService {
	return &communityService{
		repo:        repo,
		authService: authService,
	}
}

func (s *communityService) CreateCommunity(payload *CommunityPayload, token string) (*models.Community, error) {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return nil, fmt.Errorf("unauthorized")
    }
    
    community, err := s.repo.CreateCommunity(payload, userID)
    if err != nil {
        return nil, fmt.Errorf("error creating community: %w", err)
    }
    
    return community, nil
}



func (s *communityService) GetCommunity(id string) (*models.Community, error) {
	return s.repo.GetCommunity(id)
}

func (s *communityService) GetAllCommunities() ([]*models.Community, error) {
	return s.repo.GetAllCommunities()
}

func (s *communityService) JoinCommunity(token string, communityID string) error {
    userID, err := s.authService.ValidateToken(token)
    if err != nil {
        return fmt.Errorf("unauthorized")
    }

    if _, err := uuid.Parse(communityID); err != nil {
        return fmt.Errorf("invalid community UUID format: %w", err)
    }

    err = s.repo.AddUserToCommunity(userID, communityID)
    if err != nil {
        return err 
    }

    return nil
}

func (s *communityService) GetCommunityUsers(communityID string) ([]*models.User, error) {
    if _, err := uuid.Parse(communityID); err != nil {
        return nil, fmt.Errorf("invalid community UUID format: %w", err)
    }

    users, err := s.repo.GetCommunityMembers(communityID)
    if err != nil {
        return nil, err
    }
    return users, nil
}
