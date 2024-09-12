package chat

import (
	"anonymous/models"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
    db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
    return &MessageRepository{db: db}
}

func (mr *MessageRepository) CreateMessage(message *models.Message) error {
    _, err := mr.db.Exec(`
        INSERT INTO messages (id, from_user_id, to_user_id, content, sent_at)
        VALUES ($1, $2, $3, $4, $5)`,
        uuid.New().String(), message.From, message.To, message.Content, time.Now())
    if err != nil {
        return err
    }
    return nil
}

func (mr *MessageRepository) GetMessage(id string) (*models.Message, error) {
    var message models.Message
    err := mr.db.QueryRow(`
        SELECT id, from_user_id, to_user_id, content, sent_at
        FROM messages
        WHERE id = $1`, id).Scan(&message.ID, &message.From, &message.To, &message.Content, &message.SentAt)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &message, err
}
func (mr *MessageRepository) UpdateMessageContent(id string, content string) error {
    _, err := mr.db.Exec(`
        UPDATE messages
        SET content = $1
        WHERE id = $2`,
        content, id)
    if err != nil {
        return err
    }
    return nil
}

func (mr *MessageRepository) DeleteMessage(id string) error {
    _, err := mr.db.Exec(`
        DELETE FROM messages
        WHERE id = $1`, id)
    if err != nil {
        return err
    }
    return nil
}

func (mr *MessageRepository) GetMessagesBetweenUsers(user1ID, user2ID string) ([]*models.Message, error) {
    var messages []*models.Message

    query := `
        SELECT id, from_user_id, to_user_id, content, sent_at
        FROM messages
        WHERE (from_user_id = $1 AND to_user_id = $2)
           OR (from_user_id = $2 AND to_user_id = $1)
        ORDER BY sent_at`

    rows, err := mr.db.Queryx(query, user1ID, user2ID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var message models.Message
        if err := rows.StructScan(&message); err != nil {
            return nil, err
        }
        messages = append(messages, &message)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return messages, nil
}

func (mr *MessageRepository) GetMessagesByOwner(userID string) ([]models.Message, error) {
    var messages []models.Message
    err := mr.db.Select(&messages, `
        SELECT id, from_user_id, to_user_id, content, sent_at
        FROM messages
        WHERE from_user_id = $1
        ORDER BY sent_at ASC`, userID)
    if err != nil {
        return nil, err
    }
    return messages, nil
}

func (mr *MessageRepository) GetMessagesInChat(user1ID, user2ID string) ([]models.Message, error) {
    var messages []models.Message
    err := mr.db.Select(&messages, `
        SELECT id, from_user_id, to_user_id, content, sent_at,
               from_user_id = $1 as owner
        FROM messages
        WHERE (from_user_id = $1 AND to_user_id = $2)
           OR (from_user_id = $2 AND to_user_id = $1)
        ORDER BY sent_at ASC`, user1ID, user2ID)
    if err != nil {
        return nil, err
    }
    return messages, nil
}
func (mr *MessageRepository) IsMessageOwner(messageID, userID string) (bool, error) {
    var ownerID string
    err := mr.db.Get(&ownerID, `
        SELECT from_user_id
        FROM messages
        WHERE id = $1`, messageID)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }
    return ownerID == userID, nil
}

func (mr *MessageRepository) HideMessageForUser(messageID, userID string) error {
    _, err := mr.db.Exec(`
        UPDATE messages
        SET hidden_for = array_append(hidden_for, $1)
        WHERE id = $2`, userID, messageID)
    return err
}

func (mr *MessageRepository) DeleteMessageForAll(messageID string) error {
    _, err := mr.db.Exec(`
        DELETE FROM messages
        WHERE id = $1`, messageID)
    return err
}

func (mr *MessageRepository) GetConversations(userID string) ([]*models.Conversation, error) {
    query := `
        SELECT
            CASE
                WHEN from_user_id = $1 THEN to_user_id
                ELSE from_user_id
            END AS user_id,
            (SELECT username FROM users WHERE id =
                CASE
                    WHEN from_user_id = $1 THEN to_user_id
                    ELSE from_user_id
                END
            ) AS username,
            (SELECT profile_picture FROM users WHERE id =
                CASE
                    WHEN from_user_id = $1 THEN to_user_id
                    ELSE from_user_id
                END
            ) AS profile_picture,
            id AS last_message_id,
            content AS last_message_content,
            sent_at AS last_message_sent_at
        FROM messages
        WHERE from_user_id = $1 OR to_user_id = $1
        ORDER BY sent_at DESC
    `

    rows, err := mr.db.Queryx(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var conversations []*models.Conversation
    for rows.Next() {
        var conversation models.Conversation
        if err := rows.StructScan(&conversation); err != nil {
            return nil, err
        }
        conversations = append(conversations, &conversation)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return conversations, nil
}
