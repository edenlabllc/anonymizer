package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonBMap map[string]interface{}

func (p JsonBMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *JsonBMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	if err := json.Unmarshal(source, &i); err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}
