package types

import (
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	ID             uuid.UUID      `json:"id"`
	Email          NullableString `json:"email"`
	Password       NullableString `json:"password"`
	UpdatePassword NullableString `json:"update_password"`
	FirstName      NullableString `json:"first_name"`
	LastName       NullableString `json:"last_name"`
	CallSign       NullableString `json:"call_sign"`
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
	Role           NullableString `json:"role"`
	RoleRequested  NullableString `json:"role_requested,omitempty"`
	CreatedOn      NullableTime   `json:"created_on"`
	UpdatedOn      NullableTime   `json:"updated_on"`
	LastLoggedIn   NullableTime   `json:"last_logged_in,omitempty"`
}

type UserResponse struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	CallSign       string     `json:"call_sign"`
	PrimaryMDS     string     `json:"primary_mds"`
	SecondaryMDS   string     `json:"secondary_mds"`
	SSNLast4       string     `json:"ssn_last_4"`
	FlightAuthCode string     `json:"flight_auth_code"`
	IssuingUnit    string     `json:"issuing_unit"`
	UnitCharged    string     `json:"unit_charged"`
	HarmLocation   string     `json:"harm_location"`
	Status         string     `json:"status"`
	IsInstructor   bool       `json:"is_instructor"`
	IsEvaluator    bool       `json:"is_evaluator"`
	Role           string     `json:"role"`
	RoleRequested  *string    `json:"role_requested,omitempty"`
	CreatedOn      time.Time  `json:"created_on"`
	UpdatedOn      time.Time  `json:"updated_on"`
	LastLoggedIn   *time.Time `json:"last_logged_in,omitempty"`
}

type UserDbo struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	PasswordHash    string     `json:"-"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	CallSign        string     `json:"call_sign"`
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
	RoleId          uuid.UUID  `json:"role_id"`
	RoleRequestedId *uuid.UUID `json:"role_requested_id,omitempty"`
	CreatedOn       time.Time  `json:"created_on"`
	UpdatedOn       time.Time  `json:"updated_on"`
	LastLoggedIn    *time.Time `json:"last_logged_in,omitempty"`
}
