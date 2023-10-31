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
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sdbeard/go-supportlib/common/util"
)

var (
	store       *sessions.CookieStore
	sessionName string
)

/**********************************************************************************/

func InitSession(secret []byte, name string) {
	store = sessions.NewCookieStore(secret)
	sessionName = name
}

func GetSessionValue[TValue any](req *http.Request, key string) (TValue, error) {
	session, err := store.Get(req, sessionName)
	if err != nil {
		return util.GetTypeObject[TValue](), err
	}

	value, ok := session.Values[key].(TValue)
	if !ok {
		return util.GetTypeObject[TValue](), fmt.Errorf("missing session")
	}

	return value, nil
}

func SetSessionValue(req *http.Request, res http.ResponseWriter, key string, value interface{}) error {
	session, _ := store.Get(req, sessionName)

	session.Values[key] = value
	return session.Save(req, res)
}

/**********************************************************************************/
