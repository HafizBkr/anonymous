package posts

import "anonymous/validator"

type PostPayload struct {
    ID          string `json:"id"`          // Ajout de l'ID
    UserID      string `json:"user_id"`
    ContentType string `json:"content_type"`
    Content     string `json:"content"`
    Description string `json:"description"`
}

func (p *PostPayload) Validate() map[string]string {
    err := make(map[string]string)
    if validator.IsEmptyString(p.UserID) {
        err["user_id"] = "User ID is required"
    }
    if validator.IsEmptyString(p.ContentType) {
        err["content_type"] = "Content type is required"
    }
    if validator.IsEmptyString(p.Content) {
        err["content"] = "Content is required"
    }
    return err
}
