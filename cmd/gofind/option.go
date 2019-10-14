package main

import (
	"bytes"
	"encoding/json"
)

// StringOption represents an optional string
// This helps distringuish between an empty string vs unconfigured value
type StringOption struct {
	valid bool
	value string
}

func (opt *StringOption) String() string {
	return opt.value
}

// Set is called by flag.Parse
func (opt *StringOption) Set(val string) error {
	opt.valid = true
	opt.value = val
	return nil
}

func (opt *StringOption) MarshalJSON() ([]byte, error) {
	if !opt.valid {
		return json.Marshal(nil)
	}
	return json.Marshal(opt.value)
}

func (opt *StringOption) UnmarshalJSON(b []byte) error {
	if b == nil || bytes.Equal(b, []byte("null")) { // null - represents a JSON null
		return nil
	}

	if err := json.Unmarshal(b, &opt.value); err != nil {
		return err
	}

	opt.valid = true
	return nil
}

func (opt *StringOption) IsValid() bool {
	return opt.valid
}
