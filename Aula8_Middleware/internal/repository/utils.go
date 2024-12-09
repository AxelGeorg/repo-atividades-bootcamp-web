package repository

import (
	"encoding/json"
	"errors"
)

func ToBool(value interface{}) (*bool, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result bool
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, errors.New("invalid type for bool conversion")
	}
	return &result, nil
}

func ToInt(value interface{}) (int, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return 0, err
	}

	var result int
	if err := json.Unmarshal(bytes, &result); err != nil {
		return 0, errors.New("invalid type for int conversion")
	}
	return result, nil
}
