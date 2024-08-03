package communitychats

import (
	"anonymous/models"
	"anonymous/auth"
)

type CommunityChatService struct {
	repo        CommunityChatRepository
	authService *auth.AuthService
}

func NewCommunityChatService(repo CommunityChatRepository, authService *auth.AuthService) *CommunityChatService {
	return &CommunityChatService{repo: repo, authService: authService}
}

func (s *CommunityChatService) CreateMessage(chat models.CommunityChat) error {
	return s.repo.Create(chat)
}
func (s *CommunityChatService) GetMessagesByCommunityID(communityID string) ([]models.CommunityChat, error) {
	return s.repo.GetByCommunityID(communityID)
}
