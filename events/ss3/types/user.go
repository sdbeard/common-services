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

	"github.com/sdbeard/go-supportlib/common/util"
)

/**********************************************************************************/

// NewUser creates a new user object and returns the reference
func NewUser() *User {
	return &User{
		Profile: new(UserProfile),
		Claims:  make(map[string]interface{}),
		Roles:   make([]string, 0),
		Created: time.Now(),
	}
}

/***** User ***********************************************************************/

// AuthUser
type User struct {
	Profile      *UserProfile           `json:"profile,omitempty"`
	Claims       map[string]interface{} `json:"claims,omitempty"`
	Roles        []string               `json:"roles"`
	Created      time.Time              `json:"created"`
	Username     string                 `json:"username"`
	Password     string                 `json:"password"`
	Organization string                 `json:"org"`
}

/***** Marshaler interfaces *******************************************************/

// MarshalJSON is a method allowing serialization of the User
func (user User) MarshalJSON() ([]byte, error) {
	type Alias User

	return json.Marshal(&struct {
		Created int64 `json:"created"`
		Alias
	}{
		Created: user.Created.Unix(),
		Alias:   (Alias)(user),
	})
}

// UnmarshalJSON is a method implemented allowing de-serialization of the
// User
func (user *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		Created int64 `json:"created"`
		*Alias
	}{
		Alias: (*Alias)(user),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	user.Created = time.Unix(aux.Created, 0)

	return nil
}

/***** Datasource Document interface implementation *******************************/

// Item returns an object that represents the object to stored
func (user *User) Item() interface{} {
	type Alias User

	item := &struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Created int64  `json:"created"`
		*Alias
	}{
		ID:      user.Id(),
		Type:    user.Type(),
		Created: user.Created.Unix(),
		Alias:   (*Alias)(user),
	}

	return item
}

// ID returns the key/id to query and identify the event bus
func (user *User) Id() string {
	return user.Username
}

// Type returns the reflect Type representation of the current object
func (user *User) Type() string {
	return util.GetTypeName(user)
}

// IdKey returns the specific key used to query an object by ID
func (user *User) IdKey() string {
	return "id"
}

// Updates the state of the document if necessary
func (user *User) Update(username string) {
	if user.Created.Unix() < 0 {
		user.Created = time.Now()
	}
}

/***** exported functions *********************************************************/
/**********************************************************************************/
/**********************************************************************************/
