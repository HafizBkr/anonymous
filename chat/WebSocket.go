package chat

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
    "log"
    "anonymous/middleware"
    "anonymous/models"
    "github.com/gorilla/websocket"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
     "github.com/go-chi/chi/v5"
)

type HTTPMessage struct {
    To      string `json:"to"`
    Content string `json:"content"`
}

type Message struct {
    ID      string    `json:"id"`
    From    string    `json:"from"`
    To      string    `json:"to"`
    Content string    `json:"content"`
    SentAt  time.Time `json:"sent_at"`
    Owner  bool  `json:"owner"`
}

type Connection struct {
    conn *websocket.Conn
    user string
}

type Hub struct {
    connections map[string]*Connection
    broadcast   chan Message
    register    chan *Connection
    unregister  chan *Connection
    mu          sync.Mutex
}

var hub = Hub{
    connections: make(map[string]*Connection),
    broadcast:   make(chan Message),
    register:    make(chan *Connection),
    unregister:  make(chan *Connection),
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func init() {
    go runHub()
}

func runHub() {
    for {
        select {
        case conn := <-hub.register:
            hub.mu.Lock()
            hub.connections[conn.user] = conn
            hub.mu.Unlock()
        case conn := <-hub.unregister:
            hub.mu.Lock()
            if _, ok := hub.connections[conn.user]; ok {
                delete(hub.connections, conn.user)
                conn.conn.Close()
            }
            hub.mu.Unlock()
        case msg := <-hub.broadcast:
            hub.mu.Lock()
            if conn, ok := hub.connections[msg.To]; ok {
                err := conn.conn.WriteJSON(msg)
                if err != nil {
                    conn.conn.Close()
                    delete(hub.connections, msg.To)
                }
            }
            hub.mu.Unlock()
        }
    }
}

func HandleHTTPMessage(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value(middlewares.ContextKeyUser).(*models.LoggedInUser)
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var msg HTTPMessage
    err := json.NewDecoder(r.Body).Decode(&msg)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    message := Message{
        ID:      uuid.New().String(),
        From:    user.ID,
        To:      msg.To,
        Content: msg.Content,
        SentAt:  time.Now(),
        Owner:   true,
    }

    // Envoyer le message via WebSocket
    hub.broadcast <- message

    // Enregistrer le message dans la base de données
    mr := NewMessageRepository(db)
    err = mr.CreateMessage(&models.Message{
        ID:      message.ID,
        From:    message.From,
        To:      message.To,
        Content: message.Content,
        SentAt:  message.SentAt,
        Owner:     message.Owner,
    })
    if err != nil {
        log.Printf("Error saving message: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "message sent"})
}


func HandleWebSocket(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value(middlewares.ContextKeyUser).(*models.LoggedInUser)
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
        return
    }

    connection := &Connection{conn: conn, user: user.Username}
    hub.register <- connection

    // Log the successful connection
    log.Printf("WebSocket connection established for user: %s", user.Username)

    go func() {
        defer func() {
            hub.unregister <- connection
            conn.Close()
            // Log when the WebSocket connection is closed
            log.Printf("WebSocket connection closed for user: %s", user.Username)
        }()
        for {
            var msg Message
            err := conn.ReadJSON(&msg)
            if err != nil {
                break
            }
            msg.From = user.Username
            msg.ID = uuid.New().String()
            msg.SentAt = time.Now()
            hub.broadcast <- msg

            mr := NewMessageRepository(db)
            err = mr.CreateMessage(&models.Message{
                ID:      msg.ID,
                From:    msg.From,
                To:      msg.To,
                Content: msg.Content,
                SentAt:  msg.SentAt,
                Owner:     msg.Owner,
            })
            if err != nil {
                log.Printf("Error saving message: %v", err)
            }
        }
    }()
}

func UpdateMessageHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
    messageID := chi.URLParam(r, "messageID")
    user := r.Context().Value(middlewares.ContextKeyUser).(*models.LoggedInUser)
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    mr := NewMessageRepository(db)
    message, err := mr.GetMessage(messageID)
    if err != nil {
        log.Printf("Error retrieving message: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    if message == nil {
        http.Error(w, "Message not found", http.StatusNotFound)
        return
    }

    if message.From != user.ID {
        http.Error(w, "Forbidden: You are not the owner of this message", http.StatusForbidden)
        return
    }

    var updatedContent struct {
        Content string `json:"content"`
    }
    err = json.NewDecoder(r.Body).Decode(&updatedContent)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    err = mr.UpdateMessageContent(messageID, updatedContent.Content)
    if err != nil {
        log.Printf("Error updating message content: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "message content updated"})
}

func GetMessageHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        messageID := chi.URLParam(r, "messageID")
        mr := NewMessageRepository(db)
        message, err := mr.GetMessage(messageID)
        if err != nil {
            log.Printf("Error retrieving message: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        if message == nil {
            http.Error(w, "Message not found", http.StatusNotFound)
            return
        }
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(message)
    }
}


func DeleteMessageHandler(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
    messageID := chi.URLParam(r, "messageID")
    user := r.Context().Value(middlewares.ContextKeyUser).(*models.LoggedInUser)
    if user == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var request struct {
        DeleteForAll bool `json:"delete_for_all"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    mr := NewMessageRepository(db)

    if request.DeleteForAll {
        isOwner, err := mr.IsMessageOwner(messageID, user.ID)
        if err != nil {
            log.Printf("Error checking message ownership: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        if !isOwner {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        if err := mr.DeleteMessageForAll(messageID); err != nil {
            log.Printf("Error deleting message for all: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    } else {
        if err := mr.HideMessageForUser(messageID, user.ID); err != nil {
            log.Printf("Error hiding message for user: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "message deleted"})
}
func GetMessagesBetweenUsersHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := r.Context().Value(middlewares.ContextKeyUser).(*models.LoggedInUser)
        if user == nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        user1ID := chi.URLParam(r, "user1ID")
        user2ID := chi.URLParam(r, "user2ID")

        log.Printf("User1ID: %s, User2ID: %s, CurrentUser: %s", user1ID, user2ID, user.ID)
        if user.ID != user1ID && user.ID != user2ID {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        mr := NewMessageRepository(db)
        messages, err := mr.GetMessagesBetweenUsers(user1ID, user2ID)
        if err != nil {
            log.Printf("Error retrieving messages between users: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(messages)
    }
}

func GetMessagesByOwnerHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := r.Context().Value(middlewares.ContextKeyUser).(*models.LoggedInUser)
        if user == nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        mr := NewMessageRepository(db)
        messages, err := mr.GetMessagesByOwner(user.ID)
        if err != nil {
            log.Printf("Error retrieving messages for user %s: %v", user.ID, err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Répondre avec les messages en format JSON
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(messages)
    }
}

func GetMessagesInChatHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := r.Context().Value(middlewares.ContextKeyUser).(*models.LoggedInUser)
        if user == nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        user2ID := chi.URLParam(r, "user2ID")

        mr := NewMessageRepository(db)
        messages, err := mr.GetMessagesInChat(user.ID, user2ID)
        if err != nil {
            log.Printf("Error retrieving messages in chat between users %s and %s: %v", user.ID, user2ID, err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Marquer les messages avec le bon propriétaire
        for i := range messages {
            messages[i].Owner = messages[i].From == user.ID
        }

        // Répondre avec les messages en format JSON
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(messages)
    }
}