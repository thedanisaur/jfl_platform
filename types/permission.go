package types

import (
	"github.com/google/uuid"
)

type PermissionDTO struct {
	ID                  uuid.UUID `json:"id"`
	RoleID              uuid.UUID `json:"role_id"`
	Resource            string    `json:"resource"`
	Operation           string    `json:"operation"`
	Effect              string    `json:"effect"`
	ConditionType       string    `json:"cond_type"`
	ConditionExpression string    `json:"cond_expr"`
}
