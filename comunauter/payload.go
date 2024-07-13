package comunauter

import "time"

type CommunityPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Community struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *CommunityPayload) Validate() map[string]string {
	err := make(map[string]string)
	if p.Name == "" {
		err["name"] = "Name is required"
	}
	if p.Description == "" {
		err["description"] = "Description is required"
	}
	return err
}
