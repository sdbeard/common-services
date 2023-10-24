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
package types

import (
	"encoding/json"
	"time"
)

/***** JWTSecret ******************************************************************/

// JWTSecret holds the values for the JWT signing string and expiration
type JWTSecret struct {
	Key        []byte        `json:"key"`
	Expiration time.Duration `json:"expiration"`
}

/***** Marshaler interfaces *******************************************************/

// MarshalJSON is a method allowing serialization of the JWTSecret
func (secret JWTSecret) MarshalJSON() ([]byte, error) {
	type Alias JWTSecret

	return json.Marshal(&struct {
		Expiration string `json:"expiration"`
		Alias
	}{
		Expiration: secret.Expiration.String(),
		Alias:      (Alias)(secret),
	})
}

// UnmarshalJSON is a method implemented allowing de-serialization of the
// JWTSecret
func (secret *JWTSecret) UnmarshalJSON(data []byte) error {
	type Alias JWTSecret
	aux := &struct {
		Expiration string `json:"expiration"`
		*Alias
	}{
		Alias: (*Alias)(secret),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	duration, err := time.ParseDuration(aux.Expiration)
	secret.Expiration = duration

	return err
}

/**********************************************************************************/
