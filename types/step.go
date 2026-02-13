package types

import (
	"time"

	"github.com/google/uuid"
)

type StepAcknowledgementDTO struct {
	ID             uuid.UUID  `json:"id"`
	StepID         uuid.UUID  `json:"step_id"`
	UserID         uuid.UUID  `json:"user_id"`
	AcknowledgedOn *time.Time `json:"acknowledged_on,omitempty"`
}

type StepDTO struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Summary   string     `json:"summary"`
	Details   string     `json:"details"`
	Required  bool       `json:"required"`
	Status    string     `json:"status"`
	CreatedOn *time.Time `json:"created_on"`
	UpdatedOn *time.Time `json:"updated_on"`
}
