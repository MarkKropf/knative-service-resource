package concourse

import "encoding/json"

type PutParams struct {
	ImageRepository string `json:"image_repository"`
	ImageDigest     string `json:"image_digest"`
}

func ParamsFromInput(rawJson string) (*PutParams, error) {
	params := &PutParams{}

	err := json.Unmarshal([]byte(rawJson), params)
	if err != nil {
		return nil, err
	}

	return params, nil
}
