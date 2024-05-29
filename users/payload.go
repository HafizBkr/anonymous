package users

import (
	"anonymous/commons"
	"anonymous/validator"
)

type changePasswordPayload struct {
	New string `json:"new"`
	Old string `json:"old"`
}

func (p *changePasswordPayload) Validate() (err map[string]string) {
	err = map[string]string{}
	if validator.IsEmptyString(p.Old) {
		err["old"] = commons.Codes.EmptyField
		return
	}
	if validator.IsEmptyString(p.New) {
		err["new"] = commons.Codes.EmptyField
		return
	}
	return nil
}

type toggleUserStatusPayload struct {
	IDs    []string `json:"ids"`
	Active bool     `json:"active"`
}

func (p *toggleUserStatusPayload) Validate() (err map[string]string) {
	err = map[string]string{}
	for _, id := range p.IDs {
		if !validator.IsUUID(id) {
			err["ids"] = commons.Codes.InvalidField
			break
		}
	}
	if len(err) != 0 {
		return err
	}
	return nil
}
