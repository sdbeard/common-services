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
	"strings"
	"time"

	"github.com/sdbeard/go-supportlib/common/util"
)

/***** JWTSecret ******************************************************************/

// JWTSecret holds the values for the JWT signing string and expiration
type JWTSecret struct {
	Previous [][]byte      `json:"previous,omitempty"`
	Key      []byte        `json:"key"`
	Exp      time.Duration `json:"exp"`
	Name     string        `json:"name"`
}

/***** Marshaler interfaces *******************************************************/

// MarshalJSON is a method allowing serialization of the JWTSecret
func (secret JWTSecret) MarshalJSON() ([]byte, error) {
	type Alias JWTSecret

	return json.Marshal(&struct {
		Exp string `json:"exp"`
		Alias
	}{
		Exp:   secret.Exp.String(),
		Alias: (Alias)(secret),
	})
}

// UnmarshalJSON is a method implemented allowing de-serialization of the
// JWTSecret
func (secret *JWTSecret) UnmarshalJSON(data []byte) error {
	type Alias JWTSecret
	aux := &struct {
		Exp string `json:"exp"`
		*Alias
	}{
		Alias: (*Alias)(secret),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	duration, err := time.ParseDuration(aux.Exp)
	secret.Exp = duration

	return err
}

/***** Secret interface implementation ********************************************/

func (secret *JWTSecret) Secret() []byte {
	return secret.Key
}

func (secret *JWTSecret) Expiration() time.Duration {
	return secret.Exp
}

/***** Datasource Document interface implementation *******************************/

// Item returns an object that represents the object to stored
func (secret *JWTSecret) Item() interface{} {
	type Alias JWTSecret

	item := &struct {
		Exp string `json:"exp"`
		*Alias
	}{
		Exp:   secret.Exp.String(),
		Alias: (*Alias)(secret),
	}

	return item
}

// ID returns the key/id to query and identify the event bus
func (secret *JWTSecret) Id() string {
	return strings.ToLower(secret.Name)
}

// Type returns the reflect Type representation of the current object
func (secret *JWTSecret) Type() string {
	return util.GetTypeName(secret)
}

// IdKey returns the specific key used to query an object by ID
func (secret *JWTSecret) IdKey() string {
	return ""
}

// Updates the state of the document if necessary
func (secret *JWTSecret) Update(user string) {}

/**********************************************************************************/

/**********************************************************************************/
