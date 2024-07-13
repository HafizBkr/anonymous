package comunauter

import (
	"anonymous/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
	 "github.com/google/uuid"
		"strings"
)

type CommunityRepo interface {
	CreateCommunity(payload *CommunityPayload, creatorID string) (*models.Community, error)
	GetCommunity(id string) (*models.Community, error)
	GetAllCommunities() ([]*models.Community, error)
	AddUserToCommunity(userID string, communityID string) error
    GetCommunityMembers(communityID string) ([]*models.User, error)
}

type communityRepo struct {
	db *sqlx.DB
}

func NewCommunityRepo(db *sqlx.DB) CommunityRepo {
	return &communityRepo{
		db: db,
	}
}

func (r *communityRepo) CreateCommunity(payload *CommunityPayload, creatorID string) (*models.Community, error) {
    community := &models.Community{
        Name:        payload.Name,
        Description: payload.Description,
        CreatorID:   creatorID,
        CreatedAt:   time.Now(),
    }
    
    query := `
        INSERT INTO communities (id, name, description, creator_id, created_at)
        VALUES (uuid_generate_v4(), $1, $2, $3, $4)
        RETURNING id
    `
    
    err := r.db.QueryRow(query, community.Name, community.Description, community.CreatorID, community.CreatedAt).Scan(&community.ID)
    if err != nil {
        return nil, fmt.Errorf("error creating community: %w", err)
    }
    
    return community, nil
}


func (r *communityRepo) GetCommunity(id string) (*models.Community, error) {
    uuid, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid UUID format: %w", err)
    }
    
    var community models.Community
    query := "SELECT * FROM communities WHERE id = $1"
    err = r.db.Get(&community, query, uuid)
    if err != nil {
        return nil, fmt.Errorf("error getting community: %w", err)
    }
    
    return &community, nil
}



func (r *communityRepo) GetAllCommunities() ([]*models.Community, error) {
	var communities []*models.Community
	query := "SELECT * FROM communities ORDER BY created_at DESC"
	err := r.db.Select(&communities, query)
	if err != nil {
		return nil, fmt.Errorf("error getting communities: %w", err)
	}
	return communities, nil
}

func (r *communityRepo) AddUserToCommunity(userID string, communityID string) error {
    // Trim any leading or trailing whitespace
    communityID = strings.TrimSpace(communityID)
    
    // Validate communityID as UUID
    uuid, err := uuid.Parse(communityID)
    if err != nil {
        return fmt.Errorf("invalid community UUID format: %w", err)
    }

    // Check user membership
    var count int
    query := `
        SELECT COUNT(*) FROM community_members
        WHERE user_id = $1 AND community_id = $2
    `
    err = r.db.Get(&count, query, userID, uuid)
    if err != nil {
        return fmt.Errorf("error checking user membership: %w", err)
    }
    if count > 0 {
        return fmt.Errorf("user is already a member of the community")
    }

    // Add user to community
    insertQuery := `
        INSERT INTO community_members (user_id, community_id, joined_at)
        VALUES ($1, $2, $3)
    `
    _, err = r.db.Exec(insertQuery, userID, uuid, time.Now())
    if err != nil {
        return fmt.Errorf("error adding user to community: %w", err)
    }
    return nil
}

func (r *communityRepo) GetCommunityMembers(communityID string) ([]*models.User, error) {
    // Trim any leading or trailing whitespace
    communityID = strings.TrimSpace(communityID)
    
    // Validate communityID as UUID
    uuid, err := uuid.Parse(communityID)
    if err != nil {
        return nil, fmt.Errorf("invalid community UUID format: %w", err)
    }

    var members []*models.User
    query := `
        SELECT u.id, u.email, u.username, u.joined_at, u.active, u.profile_picture, u.email_verified
        FROM users u
        JOIN community_members cm ON u.id = cm.user_id
        WHERE cm.community_id = $1
    `
    if err := r.db.Select(&members, query, uuid); err != nil {
        return nil, fmt.Errorf("error retrieving community members: %w", err)
    }
    return members, nil
}
