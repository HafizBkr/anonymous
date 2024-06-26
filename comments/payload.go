package comments

import (
    "github.com/google/uuid"
    "anonymous/validator"
    "anonymous/commons"
)

type CommentPayload struct {
    PostID      string `json:"post_id"`
    UserID      string `json:"user_id"`
    ContentType string `json:"content_type"`
    Content     string `json:"content"`
}

type UpdateCommentPayload struct {
	Content string `json:"content"`
}

func (p *UpdateCommentPayload) Validate() (err map[string]string)  { 
	err = map[string]string{}
		if validator.IsEmptyString(p.Content) {
			err["label"] = commons.Codes.EmptyField
			return
		}
		return nil
}

func (p *CommentPayload) Validate() map[string]string {
    errors := make(map[string]string)
    if p.PostID == "" {
        errors["post_id"] = "Post ID is required"
    } else if !isValidUUID(p.PostID) {
        errors["post_id"] = "Post ID must be a valid UUID"
    }
    if p.UserID == "" {
        errors["user_id"] = "User ID is required"
    } else if !isValidUUID(p.UserID) {
        errors["user_id"] = "User ID must be a valid UUID"
    }
    if p.ContentType == "" {
        errors["content_type"] = "Content type is required"
    }
    if p.Content == "" {
        errors["content"] = "Content is required"
    }

    return errors
}

func isValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
}

