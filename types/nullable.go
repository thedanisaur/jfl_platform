// These types allow for a user to either explicitly set certain values to null in a json object
// or omit the values from the object so that we can selectively update only the provided values.
package types

import (
	"encoding/json"
	"time"
)

type NullableString struct {
	Value *string
	Set   bool
}

func (ns *NullableString) UnmarshalJSON(data []byte) error {
	ns.Set = true
	// Check to see if we're trying to set this value to null.
	if string(data) == "null" {
		ns.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	ns.Value = &s
	return nil
}

type NullableBool struct {
	Value *bool
	Set   bool
}

func (nb *NullableBool) UnmarshalJSON(data []byte) error {
	nb.Set = true
	// Check to see if we're trying to set this value to null.
	if string(data) == "null" {
		nb.Value = nil
		return nil
	}
	var b bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	nb.Value = &b
	return nil
}

type NullableBytes struct {
	Value *[]byte
	Set   bool
}

func (nb *NullableBytes) UnmarshalJSON(data []byte) error {
	nb.Set = true
	// Check to see if we're trying to set this value to null.
	if string(data) == "null" {
		nb.Value = nil
		return nil
	}
	var b []byte
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	nb.Value = &b
	return nil
}

type NullableTime struct {
	Value *time.Time
	Set   bool
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	nt.Set = true
	// Check to see if we're trying to set this value to null.
	if string(data) == "null" {
		nt.Value = nil
		return nil
	}
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	nt.Value = &t
	return nil
}
