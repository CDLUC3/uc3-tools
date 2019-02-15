package storagenodes

import (
	"errors"
	props "github.com/magiconair/properties"
)

type CloudCredentials struct {
	Key    string
	Secret string
}

func LoadCredentials(svcProps *props.Properties) (*CloudCredentials, error) {
	key := svcProps.GetString("accessKey", "")
	if key == "" {
		key = svcProps.GetString("access_key", "")
	}
	secret := svcProps.GetString("secretKey", "")
	if secret == "" {
		secret = svcProps.GetString("secret_key", "")
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