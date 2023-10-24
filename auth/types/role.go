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
/***** Role ***********************************************************************/

// Role defines a Role that a User holds as part of an RBAC system
type Role struct {
	Created     time.Time `json:"created"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
}

/***** Marshaler interface implmentation ******************************************/

// MarshalJSON is a method allowing serialization of the Role
func (role Role) MarshalJSON() ([]byte, error) {
	type Alias Role

	return json.Marshal(&struct {
		Created int64 `json:"created"`
		Alias
	}{
		Created: role.Created.Unix(),
		Alias:   (Alias)(role),
	})
}

// UnmarshalJSON is a method implemented allowing de-serialization of the
// Role
func (role *Role) UnmarshalJSON(data []byte) error {
	type Alias Role
	aux := &struct {
		Created int64 `json:"created"`
		*Alias
	}{
		Alias: (*Alias)(role),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	role.Created = time.Unix(aux.Created, 0)

	return nil
}

/***** Datasource Document interface implementation *******************************/

// Item returns an object that represents the object to stored
func (role *Role) Item() interface{} {
	type Alias Role

	item := &struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Created int64  `json:"created"`
		*Alias
	}{
		ID:      role.Id(),
		Type:    role.Type(),
		Created: role.Created.Unix(),
		Alias:   (*Alias)(role),
	}

	return item
}

// ID returns the key/id to query and identify the event bus
func (role *Role) Id() string {
	return role.Name
}

// Type returns the reflect Type representation of the current object
func (role *Role) Type() string {
	return util.GetTypeName(role)
}

// IdKey returns the specific key used to query an object by ID
func (role *Role) IdKey() string {
	return "id"
}

// Updates the state of the document if necessary
func (role *Role) Update(user string) {
	if role.Created.Unix() < 0 {
		role.Created = time.Now()
	}
}

/**********************************************************************************/
/**********************************************************************************/
