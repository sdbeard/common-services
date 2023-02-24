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
package main

import (
	"encoding/json"
	"time"

	"github.com/sdbeard/go-supportlib/common/util"
)

/***** User ***********************************************************************/

type User struct {
	Created   time.Time `json:"created" dynamodbav:"created"`
	Updated   time.Time `json:"updated" dynamodbav:"updated"`
	Name      string    `json:"name" dynamodbav:"name"`
	Email     string    `json:"email" dynamodbav:"email" ddb:"condition=attribute_not_exists(email)|hashkey"`
	Pwd       string    `json:"password" dynamodbav:"password"`
	Role      string    `json:"role" dynamodbav:"role"`
	IsEnabled bool      `json:"enabled" dynamodbav:"enabled"`
}

/***** Marshaler interface implementations ****************************************/

// MarshalJSON marshals the Organization to a JSON string
func (user User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Created int64 `json:"created"`
		Updated int64 `json:"updated"`
		Alias
	}{
		Created: user.Created.Unix(),
		Updated: user.Updated.Unix(),
		Alias:   (Alias)(user),
	})
}

// UnmarshalJSON unmarshals JSON string to an Organization object
func (user *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		Created int64 `json:"created"`
		Updated int64 `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(user),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	user.Created = time.Unix(aux.Created, 0)
	user.Updated = time.Unix(aux.Updated, 0)

	return nil
}

/***** Datasource Document interface implementation *******************************/

// ID returns the key/id to query and identify the event bus
func (user *User) Id() string {
	return user.Email
}

// Item returns an anonymous struct to be saved
func (user *User) Item() interface{} {
	/*
		type Alias User

		ddbUser := &struct {
			PK   string `dynamodbav:"PK" ddb:"condition=attribute_not_exists(PK)|hashkey=PK`
			Type string `json:"type" dynamodbav:"type"`
			*Alias
		}{
			PK:    user.Id(),
			Type:  user.Type(),
			Alias: (*Alias)(user),
		}

		return *ddbUser
	*/
	return user
}

func (user *User) Type() string {
	return util.GetTypeName(user)
}

// Update makes sure the current Organization has it's variable values updated for
// tracking and storage
func (user *User) Update(userName string) {
	now := time.Now()

	if user.Created.Unix() <= 0 {
		user.Created = now
	}

	user.Updated = now
}

func (user *User) UserId() string {
	return user.Email
}

func (user *User) Password() string {
	return user.Pwd
}

func (user *User) Enabled() bool {
	return user.IsEnabled
}

func (user *User) SetPassword(password string) {
	user.Pwd = password
}

func (user *User) GetClaims() map[string]interface{} {
	return map[string]interface{}{
		"user": user.UserId(),
		"role": user.Role,
	}
}

/**********************************************************************************/
