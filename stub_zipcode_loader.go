package main

import "context"

// NewStubZipcodeLoader creates a new stubbed zipcode loader that
// always returns the same zipcode
func NewStubZipcodeLoader(zip string) *StubZipcodeLoader {
	return &StubZipcodeLoader{zipcode: zip}
}

// StubZipcodeLoader is a stub zipcode loader that returns the same
// zipcode for every request
type StubZipcodeLoader struct {
	zipcode string
}

func (s *StubZipcodeLoader) DeviceZip(ctx context.Context, deviceID, apiHost, apiAccessToken string) (string, error) {
	return s.zipcode, nil
}
