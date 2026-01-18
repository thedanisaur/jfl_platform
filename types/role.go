package types

import (
	"github.com/google/uuid"
)

type RoleDto struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayname"`
}

type RoleDbo struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayname"`
}
