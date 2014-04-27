package chef

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Services map[string]interface{}  `json:"services"`
}

func Unmarshal(data []byte) (*Config, error) {
	var config *Config
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}


func UnmarshalFile(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := Unmarshal(file)
	return config, err
}

func (config *Config) Marshal() ([]byte, error) {
	json, err := json.MarshalIndent(config, "", "  ")
	return json, err
}