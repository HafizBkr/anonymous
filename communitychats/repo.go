package communitychats

import (
	"github.com/jmoiron/sqlx"
	"errors"
	"anonymous/models"
	"github.com/google/uuid"
)

type CommunityChatRepository interface {
	Create(chat models.CommunityChat) error
	GetByCommunityID(communityID string) ([]models.CommunityChat, error)
}

type communityChatRepo struct {
	db *sqlx.DB
}

func NewCommunityChatRepo(db *sqlx.DB) CommunityChatRepository {
	return &communityChatRepo{db: db}
}

func (repo *communityChatRepo) Create(chat models.CommunityChat) error {
    chat.ID = uuid.New().String() // Génération d'un nouvel UUID
    query := `INSERT INTO community_chats (id, community_id, user_id, message, created_at) VALUES ($1, $2, $3, $4, $5)`

    _, err := repo.db.Exec(query, chat.ID, chat.CommunityID, chat.UserID, chat.Message, chat.CreatedAt)
    if err != nil {
        return errors.New("failed to insert community chat message: " + err.Error())
    }

    return nil
}

func (repo *communityChatRepo) GetByCommunityID(communityID string) ([]models.CommunityChat, error) {
	query := `SELECT id, community_id, user_id, message, created_at FROM community_chats WHERE community_id = $1`
	rows, err := repo.db.Query(query, communityID)
	if err != nil {
		return nil, errors.New("failed to get community chat messages: " + err.Error())
	}
	defer rows.Close()

	var chats []models.CommunityChat
	for rows.Next() {
		var chat models.CommunityChat
		if err := rows.Scan(&chat.ID, &chat.CommunityID, &chat.UserID, &chat.Message, &chat.CreatedAt); err != nil {
			return nil, errors.New("failed to scan community chat message: " + err.Error())
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
