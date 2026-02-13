package types

import (
	"time"

	"github.com/google/uuid"
)

type FlightLogAircrewDTO struct {
	ID                          uuid.UUID  `json:"id"`
	FlightLogID                 uuid.UUID  `json:"flight_log_id"`
	UserID                      uuid.UUID  `json:"user_id"`
	FlyingOrigin                string     `json:"flying_origin"`
	FlightAuthCode              string     `json:"flight_auth_code"`
	TimePrimary                 float64    `json:"time_primary"`
	TimeSecondary               float64    `json:"time_secondary"`
	TimeInstructor              float64    `json:"time_instructor"`
	TimeEvaluator               float64    `json:"time_evaluator"`
	TimeOther                   float64    `json:"time_other"`
	TotalAircrewDurationDecimal float64    `json:"total_aircrew_duration_decimal"`
	TotalAircrewSorties         int64      `json:"total_aircrew_sorties"`
	CondNightTime               float64    `json:"cond_night_time"`
	CondInstrumentTime          float64    `json:"cond_instrument_time"`
	CondSimInstrumentTime       float64    `json:"cond_sim_instrument_time"`
	CondNvgTime                 float64    `json:"cond_nvg_time"`
	CondCombatTime              float64    `json:"cond_combat_time"`
	CondCombatSortie            int64      `json:"cond_combat_sortie"`
	CondCombatSupportTime       float64    `json:"cond_combat_support_time"`
	CondCombatSupportSortie     int64      `json:"cond_combat_support_sortie"`
	AircrewRoleType             string     `json:"aircrew_role_type"`
	CreatedOn                   *time.Time `json:"created_on"`
	UpdatedOn                   *time.Time `json:"updated_on"`
}

type FlightLogCommentDTO struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	FlightLogID uuid.UUID  `json:"flight_log_id,omitempty"`
	UserID      uuid.UUID  `json:"user_id"`
	RoleName    string     `json:"role_name"`
	Comment     string     `json:"comment"`
	CreatedOn   *time.Time `json:"created_on"`
	UpdatedOn   *time.Time `json:"updated_on"`
}

type FlightLogDTO struct {
	ID                         uuid.UUID             `json:"id"`
	UserID                     uuid.UUID             `json:"user_id"`
	MDS                        string                `json:"mds"`
	FlightLogDate              *time.Time            `json:"flight_log_date"`
	SerialNumber               string                `json:"serial_number"`
	UnitCharged                string                `json:"unit_charged"`
	HarmLocation               string                `json:"harm_location"`
	FlightAuthorization        string                `json:"flight_authorization"`
	IssuingUnit                string                `json:"issuing_unit"`
	IsTrainingFlight           bool                  `json:"is_training_flight"`
	IsTrainingOnly             bool                  `json:"is_training_only"`
	TotalFlightDecimalTime     float64               `json:"total_flight_decimal_time"`
	SchedulerSignatureID       *uuid.UUID            `json:"scheduler_signature_id,omitempty"`
	SarmSignatureID            *uuid.UUID            `json:"sarm_signature_id,omitempty"`
	InstructorSignatureID      *uuid.UUID            `json:"instructor_signature_id,omitempty"`
	StudentSignatureID         *uuid.UUID            `json:"student_signature_id,omitempty"`
	TrainingOfficerSignatureID *uuid.UUID            `json:"training_officer_signature_id,omitempty"`
	Type                       string                `json:"type,omitempty"`
	Remarks                    string                `json:"remarks,omitempty"`
	Missions                   []FlightLogMissionDTO `json:"missions,omitempty"`
	Aircrew                    []FlightLogAircrewDTO `json:"aircrew,omitempty"`
	Comments                   []FlightLogCommentDTO `json:"comments,omitempty"`
}

type FlightLogMissionDTO struct {
	ID               uuid.UUID  `json:"id"`
	FlightLogID      uuid.UUID  `json:"flight_log_id"`
	MissionNumber    string     `json:"mission_number,omitempty"`
	MissionSymbol    string     `json:"mission_symbol"`
	MissionFrom      string     `json:"mission_from"`
	MissionTo        string     `json:"mission_to"`
	TakeoffTime      *time.Time `json:"takeoff_time,omitempty"`
	LandTime         *time.Time `json:"land_time,omitempty"`
	TotalTimeDecimal float64    `json:"total_time_decimal"`
	TotalTimeDisplay string     `json:"total_time_display,omitempty"`
	TouchAndGos      int64      `json:"touch_and_gos"`
	FullStops        int64      `json:"full_stops"`
	TotalLandings    int64      `json:"total_landings"`
	Sorties          int64      `json:"sorties"`
	CreatedOn        *time.Time `json:"created_on"`
	UpdatedOn        *time.Time `json:"updated_on"`
}

type TemplateFlightLogDTO struct {
	FlightLogDTO
	Name string `json:"name"`
}
