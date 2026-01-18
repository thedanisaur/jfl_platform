package types

import (
	"time"

	"github.com/google/uuid"
)

type UserDto struct {
	ID             uuid.UUID      `json:"id"`
	Email          NullableString `json:"email"`
	Password       NullableString `json:"password"`
	UpdatePassword NullableString `json:"update_password"`
	Firstname      NullableString `json:"firstname"`
	Lastname       NullableString `json:"lastname"`
	Callsign       NullableString `json:"callsign"`
	PrimaryMDS     NullableString `json:"primary_mds"`
	SecondaryMDS   NullableString `json:"secondary_mds"`
	SSNLast4       NullableString `json:"ssn_last_4"`
	FlightAuthCode NullableString `json:"flight_auth_code"`
	IssuingUnit    NullableString `json:"issuing_unit"`
	UnitCharged    NullableString `json:"unit_charged"`
	HarmLocation   NullableString `json:"harm_location"`
	Status         NullableString `json:"status"`
	IsInstructor   NullableBool   `json:"is_instructor"`
	IsEvaluator    NullableBool   `json:"is_evaluator"`
	RoleID         NullableBytes  `json:"role_id"`
	RoleRequested  NullableBytes  `json:"role_requested_id,omitempty"`
	CreatedOn      NullableTime   `json:"created_on"`
	UpdatedOn      NullableTime   `json:"updated_on"`
	LastLoggedIn   NullableTime   `json:"last_logged_in,omitempty"`
}

type UserDbo struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	PasswordHash    string     `json:"-"`
	Firstname       string     `json:"firstname"`
	Lastname        string     `json:"lastname"`
	Callsign        string     `json:"callsign"`
	PrimaryMDS      string     `json:"primary_mds"`
	SecondaryMDS    string     `json:"secondary_mds"`
	SSNLast4        string     `json:"ssn_last_4"`
	FlightAuthCode  string     `json:"flight_auth_code"`
	IssuingUnit     string     `json:"issuing_unit"`
	UnitCharged     string     `json:"unit_charged"`
	HarmLocation    string     `json:"harm_location"`
	Status          string     `json:"status"`
	IsInstructor    bool       `json:"is_instructor"`
	IsEvaluator     bool       `json:"is_evaluator"`
	RoleID          uuid.UUID  `json:"role_id"`
	RoleRequestedId *uuid.UUID `json:"role_requested_id,omitempty"`
	CreatedOn       time.Time  `json:"created_on"`
	UpdatedOn       time.Time  `json:"updated_on"`
	LastLoggedIn    *time.Time `json:"last_logged_in,omitempty"`
}
