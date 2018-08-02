package concourse

import "encoding/json"

type Params struct {
	ImageRepository string `json:"image_repository"`
	ImageDigest     string `json:"image_digest"`
}

func ParamsFromInput(rawJson string) (*Params, error) {
	params := &Params{}

	err := json.Unmarshal([]byte(rawJson), params)
	if err != nil {
		return nil, err
	}

	return params, nil
}
