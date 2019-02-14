package storage

import (
	"errors"
	props "github.com/magiconair/properties"
)

type CloudCredentials struct {
	Key    string
	Secret string
}

func LoadCredentials(nodeProps *props.Properties) (*CloudCredentials, error) {
	key := nodeProps.GetString("accessKey", "")
	if key == "" {
		key = nodeProps.GetString("access_key", "")
	}
	secret := nodeProps.GetString("secretKey", "")
	if secret == "" {
		secret = nodeProps.GetString("secret_key", "")
	}

	if key == "" {
		if secret == "" {
			return nil, nil
		}
		return nil, errors.New("can't provide secret without access key")
	}
	if secret == "" {
		return nil, errors.New("can't provide access key without secret")
	}
	return &CloudCredentials{Key: key, Secret: secret}, nil
}