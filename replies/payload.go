package replies
import (
    "github.com/google/uuid"
    "anonymous/validator"
    "anonymous/commons"
)

type CommentReplyPayload struct {
    CommentID   string `json:"comment_id"`
    UserID      string `json:"user_id"`
    ContentType string `json:"content_type"`
    Content     string `json:"content"`
}

type UpdateCommentReplyPayload struct {
    Content string `json:"content"`
}

func (p *UpdateCommentReplyPayload) Validate() map[string]string {
    errors := map[string]string{}
    if validator.IsEmptyString(p.Content) {
        errors["content"] = commons.Codes.EmptyField
    }
    return errors
}

func (p *CommentReplyPayload) Validate() map[string]string {
    errors := make(map[string]string)
    if p.CommentID == "" {
        errors["comment_id"] = "Comment ID is required"
    } else if !isValidUUID(p.CommentID) {
        errors["comment_id"] = "Comment ID must be a valid UUID"
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
