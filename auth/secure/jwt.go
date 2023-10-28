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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sdbeard/common-services/auth/types"
)

/***** exported functions *********************************************************/

// GenerateJWT created the
func GenerateJWT(secret []byte, user *types.User) (string, error) {
	// Create the claims for the user token
	claims := jwt.MapClaims{
		"sub":        user.Id(),
		"roles":      user.Roles,
		"authorized": true,
		"exp":        time.Now().Add(1 * time.Minute).Unix(),
	}

	for key, value := range user.Claims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// GenerateJWT created the
func GenerateRefreshJWT(secret []byte, user *types.User) (string, error) {
	// Create the claims for the user token
	claims := jwt.MapClaims{
		"sub":        user.Id(),
		"roles":      user.Roles,
		"authorized": true,
		"exp":        time.Now().Add(1 * time.Hour * 24).Unix(),
	}

	for key, value := range user.Claims {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

/**********************************************************************************/
