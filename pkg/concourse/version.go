package concourse

import "encoding/json"

type Version struct {
	ConfigurationGeneration string `json:"configuration_generation"`
}

func VersionFromInput(rawJson string) (*Version, error) {
	version := &Version{}

	err := json.Unmarshal([]byte(rawJson), version)
	if err != nil {
		return nil, err
	}

	return version, nil
}