package db

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"user_service/types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetPassword(username string) (string, error) {
	txid := uuid.New()
	log.Printf("GetPassword | %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return "", errors.New("failed to connect to DB")
	}
	query := `
		SELECT password_hash
		FROM users
		WHERE email = LOWER(?)
	`
	row := database.QueryRow(query, username)
	var password_hash string
	err = row.Scan(&password_hash)
	if err != nil {
		log.Printf("Invalid username: %s\n", err.Error())
		return "", errors.New("failed to connect to DB")
	}
	return password_hash, nil
}

func GetUser(user_id uuid.UUID) (types.UserDbo, error) {
	txid := uuid.New()
	log.Printf("GetUser | %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return types.UserDbo{}, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) id
			, email
			, password_hash
			, firstname
			, lastname
			, callsign
			, primary_mds
			, secondary_mds
			, ssn_last_4
			, flight_auth_code
			, issuing_unit
			, unit_charged
			, harm_location
			, status
			, is_instructor
			, is_evaluator
			, role_id
			, role_requested_id
			, created_on
			, updated_on
			, last_logged_in
		FROM users
		WHERE id = UUID_TO_BIN(?)
	`
	row := database.QueryRow(query, user_id)
	var user types.UserDbo
	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Firstname,
		&user.Lastname,
		&user.Callsign,
		&user.PrimaryMDS,
		&user.SecondaryMDS,
		&user.SSNLast4,
		&user.FlightAuthCode,
		&user.IssuingUnit,
		&user.UnitCharged,
		&user.HarmLocation,
		&user.Status,
		&user.IsInstructor,
		&user.IsEvaluator,
		&user.RoleID,
		&user.RoleRequestedId,
		&user.CreatedOn,
		&user.UpdatedOn,
		&user.LastLoggedIn,
	)
	if err != nil {
		log.Printf("Failed to retrieve user:\n%s\n", err.Error())
		return types.UserDbo{}, errors.New("failed to retrieve user")
	}
	return user, nil
}

func GetUserId(email string) (uuid.UUID, error) {
	txid := uuid.New()
	log.Printf("GetUser | %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return uuid.Nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) id
		FROM users
		WHERE email = ?
	`
	row := database.QueryRow(query, email)
	var id uuid.UUID
	err = row.Scan(&id)
	if err != nil {
		log.Printf("Failed to retrieve user id:\n%s\n", err.Error())
		return uuid.Nil, errors.New("failed to retrieve user id")
	}
	return id, nil
}

func GetUsers() ([]types.UserDbo, error) {
	txid := uuid.New()
	log.Printf("GetUsers | %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return []types.UserDbo{}, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) id
			, email
			, password_hash
			, firstname
			, lastname
			, callsign
			, primary_mds
			, secondary_mds
			, ssn_last_4
			, flight_auth_code
			, issuing_unit
			, unit_charged
			, harm_location
			, status
			, is_instructor
			, is_evaluator
			, role_id
			, role_requested_id
			, created_on
			, updated_on
			, last_logged_in
		FROM users
	`
	rows, err := database.Query(query)
	if err != nil {
		log.Printf("Failed to query DB:\n%s\n", err.Error())
		return []types.UserDbo{}, errors.New("failed to connect to DB")
	}

	var users []types.UserDbo
	for rows.Next() {
		var user types.UserDbo
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Firstname,
			&user.Lastname,
			&user.Callsign,
			&user.PrimaryMDS,
			&user.SecondaryMDS,
			&user.SSNLast4,
			&user.FlightAuthCode,
			&user.IssuingUnit,
			&user.UnitCharged,
			&user.HarmLocation,
			&user.Status,
			&user.IsInstructor,
			&user.IsEvaluator,
			&user.RoleID,
			&user.RoleRequestedId,
			&user.CreatedOn,
			&user.UpdatedOn,
			&user.LastLoggedIn,
		)
		if err != nil {
			log.Printf("Failed to scan row\n%s\n", err.Error())
			continue
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		log.Println("Error scanning rows")
		return []types.UserDbo{}, errors.New("failed to connect to DB")
	}
	return users, nil
}

func InsertUser(hashed_password string, user types.UserDbo) (int64, error) {
	txid := uuid.New()
	log.Printf("InsertUser | %s\n", txid.String())
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return -1, errors.New(err_string)
	}
	log.Printf("%v\n", user)
	query := `
		INSERT INTO users
		(
			email
			, password_hash
			, firstname
			, lastname
			, callsign
			, primary_mds
			, secondary_mds
			, ssn_last_4
			, flight_auth_code
			, issuing_unit
			, unit_charged
			, harm_location
			, status
			, is_instructor
			, is_evaluator
			, role_id
			, role_requested_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UUID_TO_BIN(?), UUID_TO_BIN(?))
	`
	result, err := database.Exec(query,
		user.Email,
		hashed_password,
		user.Firstname,
		user.Lastname,
		user.Callsign,
		user.PrimaryMDS,
		user.SecondaryMDS,
		user.SSNLast4,
		user.FlightAuthCode,
		user.IssuingUnit,
		user.UnitCharged,
		user.HarmLocation,
		user.Status,
		user.IsInstructor,
		user.IsEvaluator,
		user.RoleID,
		user.RoleRequestedId,
	)
	if err != nil {
		log.Printf("Failed user insert\n%s\n", err.Error())
		return -1, errors.New(err_string)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
		return -1, errors.New(err_string)
	}
	return id, nil
}

// TODO [drd] fix this crap
func UpdateUser(current_user types.UserDbo, update_user types.UserDto) (int64, error) {
	var set_clauses []string
	var arguments []interface{}

	set_clauses, arguments = AddNullableString("email", update_user.Email, set_clauses, arguments)
	if update_user.UpdatePassword.Set {
		// Compare the password sent to the password we expect.
		err := bcrypt.CompareHashAndPassword([]byte(current_user.PasswordHash), []byte(*update_user.Password.Value))
		if err != nil {
			return -1, errors.New("invalid password")
		}
		// Now hash the new password
		hashed_password, err := bcrypt.GenerateFromPassword([]byte(*update_user.UpdatePassword.Value), 12)
		if err != nil {
			return -1, errors.New("failed to hash password")
		}
		update_hashed_password := string(hashed_password)
		// Make sure to override the old password
		update_user.Password = types.NullableString{Set: true, Value: &update_hashed_password}
		set_clauses, arguments = AddNullableString("password_hash", update_user.Password, set_clauses, arguments)
	}
	set_clauses, arguments = AddNullableString("firstname", update_user.Firstname, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("lastname", update_user.Lastname, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("callsign", update_user.Callsign, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("primary_mds", update_user.PrimaryMDS, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("secondary_mds", update_user.SecondaryMDS, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("ssn_last_4", update_user.SSNLast4, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("flight_auth_code", update_user.FlightAuthCode, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("issuing_unit", update_user.IssuingUnit, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("unit_charged", update_user.UnitCharged, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("harm_location", update_user.HarmLocation, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("status", update_user.Status, set_clauses, arguments)
	set_clauses, arguments = AddNullableBool("is_instructor", update_user.IsInstructor, set_clauses, arguments)
	set_clauses, arguments = AddNullableBool("is_evaluator", update_user.IsEvaluator, set_clauses, arguments)
	set_clauses, arguments = AddNullableBytes("role_id", update_user.RoleID, set_clauses, arguments)
	set_clauses, arguments = AddNullableBytes("role_requested_id", update_user.RoleRequested, set_clauses, arguments)
	set_clauses, arguments = AddNullableTime("created_on", update_user.CreatedOn, set_clauses, arguments)
	set_clauses, arguments = AddNullableTime("updated_on", update_user.UpdatedOn, set_clauses, arguments)
	set_clauses, arguments = AddNullableTime("last_logged_in", update_user.LastLoggedIn, set_clauses, arguments)

	if len(set_clauses) == 0 {
		return -1, fmt.Errorf("no fields to update")
	}

	query := bytes.Buffer{}
	query.WriteString("UPDATE users SET ")
	query.WriteString(set_clauses[0])
	for _, clause := range set_clauses[1:] {
		query.WriteString(", ")
		query.WriteString(clause)
	}

	query.WriteString(" WHERE id = UUID_TO_BIN(?)")
	arguments = append(arguments, update_user.ID)

	fmt.Println(query.String())
	fmt.Println(arguments)

	result, err := database.Exec(query.String(), arguments...)
	if err != nil {
		log.Printf("Failed user update\n%s\n", err.Error())
		return -1, errors.New("database error")
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed retrieve count\n%s\n", err.Error())
		return -1, errors.New("database error")
	}
	return count, nil
}
