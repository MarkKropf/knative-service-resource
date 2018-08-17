package config

type PutParams struct {
	ImageRepository string `json:"image_repository"`
	ImageDigest     string `json:"image_digest"`
}
