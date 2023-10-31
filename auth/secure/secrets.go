// *********************************************************************************
// The MIT License (MIT)
//
// # Copyright (c) 2023 Sean Beard
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in the
// Software without restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the
// Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN
// AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// *********************************************************************************
package secure

import (
	"context"

	"github.com/sdbeard/common-services/auth/conf"
	"github.com/sdbeard/go-supportlib/secure/secrets"
	"github.com/sdbeard/go-supportlib/secure/secrets/factory"
	sectypes "github.com/sdbeard/go-supportlib/secure/types"
)

var (
	secretMap = make(map[string]*sectypes.SimpleSecret)
)

/**** exported functions **********************************************************/

func LoadSecrets() error {
	return load()
}

func GetSecret(name string) (*sectypes.SimpleSecret, error) {
	secret, ok := secretMap[name]
	if !ok {
		return getSecret(name)
	}

	return secret, nil
}

func AddNewSecret(name string, size, expiry int64) (*sectypes.SimpleSecret, error) {
	if secret, ok := secretMap[name]; ok {
		return secret, nil
	}

	secret, err := createSecret(name, size, expiry)
	if err != nil {
		return nil, err
	}

	secretMap[name] = secret

	return secret, nil
}

/**********************************************************************************/

/*
func getAllSecrets(createIfMissing bool) error {
	manager, err := getSecretsManager()
	if err != nil {
		return err
	}

	secret, err := getSecret("jwtsecret")
	if err != nil {
		return err
	}
	if secret == nil && createIfMissing {
	}

	return nil
}
*/

func getSecret(name string) (*sectypes.SimpleSecret, error) {
	manager, err := getSecretsManager()
	if err != nil {
		return nil, err
	}

	secrets, err := manager.Retrieve(
		manager.Retrieve.WithSecretName(name),
	)

	return secrets[0], err
}

func createSecret(name string, size, expiry int64) (*sectypes.SimpleSecret, error) {
	manager, err := getSecretsManager()
	if err != nil {
		return nil, err
	}

	secret := sectypes.NewSimpleSecret(name, size, expiry)
	err = manager.Create(
		secret,
		manager.Create.WithContext(context.TODO()),
		manager.Create.WithAllowUpdate(false),
	)

	return secret, err
}

func getSecretsManager() (*secrets.Manager[*sectypes.SimpleSecret], error) {
	return factory.SecretsManagerFactory[*sectypes.SimpleSecret](conf.Get().SecretsConf)
}

/**********************************************************************************/

// func getAuthServiceSecrets() (*sectypes.SimpleSecret, *sectypes.SimpleSecret, *sectypes.SimpleSecret, error) {
func load() error {
	manager, err := getSecretsManager()
	if err != nil {
		return err
	}

	foundSecrets, err := manager.Retrieve(
		manager.Retrieve.WithRetrieveAll(),
	)
	if err != nil {
		return err
	}

	for _, secret := range foundSecrets {
		secretMap[secret.Id()] = secret
	}

	return err

	/*
		// Retrieve the jwt secret and refresh keys
		jwtSecret, err := getSecret("jwtsecretkey")
		if err != nil {
			return nil, nil, nil, err
		}
		if jwtSecret == nil {
			jwtSecret, _ = AddNewSecret("jwtsecretkey", 16, 60)
		}
	*/

	/*
		jwtRefreshSecret, err := getSecret("jwtrefreshsecretkey", true)
		if err != nil {
			return nil, nil, nil, err
		}

		sessionSecret, err := getSecret("sessionkey", true)
		if err != nil {
			return nil, nil, nil, err
		}
	*/

	//return jwtSecret, jwtRefreshSecret, sessionSecret, nil
	//return jwtSecret, nil, nil, nil
}
