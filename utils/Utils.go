package utils

import (
	//"encoding/json"
	//"database/sql/driver"
	"errors"
	//"fmt"
	"database/sql/driver"
	"bytes"
)

type JSONRaw []byte


func (j JSONRaw) Value() (driver.Value, error) {
	if j ==nil {
		//      log.Trace("returning null")
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONRaw) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		errors.New("Scan source was not string")
	}
	// I think I need to make a copy of the bytes.
	// It seems the byte slice passed in is re-used
	*j = append((*j)[0:0], s...)

	return nil
}


// MarshalJSON returns *m as the JSON encoding of m.
func (m JSONRaw) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *JSONRaw) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (j JSONRaw) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSONRaw) Equals(j1 JSONRaw) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

/*
func (m *JSONRaw) MarshalJSON() ([]byte, error) {
	return *m, nil
}

func (m *JSONRaw) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
} */
