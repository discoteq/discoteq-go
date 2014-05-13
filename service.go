package main

import (
	"encoding/json"
)

type Service struct {
	Name  string
	Hosts []map[string]string
}

func (service *Service) Marshall() ([]byte, error) {
	json, err := json.MarshalIndent(service, "", "  ")
	return json, err
}
