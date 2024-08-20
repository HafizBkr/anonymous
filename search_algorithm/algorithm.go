package searchalgorithm

import (
	"anonymous/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SearchService interface {
	Search(query string, limit int, offset int) (*SearchResults, error)
}

type searchService struct {
	db *sqlx.DB
}

func NewSearchService(db *sqlx.DB) SearchService {
	return &searchService{db: db}
}

type SearchResults struct {
	Users       []*models.User      `json:"users"`
	Posts       []*models.Post      `json:"posts"`
	Communities []*models.Community `json:"communities"`
}

func (s *searchService) Search(query string, limit int, offset int) (*SearchResults, error) {
	var results SearchResults

	// Recherche des utilisateurs
	usersQuery := `SELECT id, username, email, profile_picture FROM users WHERE username ILIKE $1`
	usersRows, err := s.db.Queryx(usersQuery, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("error searching users: %w", err)
	}
	defer usersRows.Close()

	for usersRows.Next() {
		var user models.User
		if err := usersRows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		results.Users = append(results.Users, &user)
	}

	// Recherche des posts
	postsQuery := `
        SELECT
            p.id,
            p.user_id,
            p.content,
            p.created_at,
            p.content_type,
            u.username
        FROM
            posts p
        JOIN
            users u ON p.user_id = u.id
        WHERE
            p.content ILIKE $1
        ORDER BY
            p.created_at DESC
        LIMIT
            $2 OFFSET $3
    `
	postsRows, err := s.db.Queryx(postsQuery, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error searching posts: %w", err)
	}
	defer postsRows.Close()

	for postsRows.Next() {
		var post models.Post
		if err := postsRows.StructScan(&post); err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		results.Posts = append(results.Posts, &post)
	}

	// Recherche des communautés
	communitiesQuery := `SELECT id, name, description, creator_id, created_at FROM communities WHERE name ILIKE $1 OR description ILIKE $1`
	communitiesRows, err := s.db.Queryx(communitiesQuery, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("error searching communities: %w", err)
	}
	defer communitiesRows.Close()

	for communitiesRows.Next() {
		var community models.Community
		if err := communitiesRows.StructScan(&community); err != nil {
			return nil, fmt.Errorf("error scanning community: %w", err)
		}
		results.Communities = append(results.Communities, &community)
	}

	return &results, nil
}
