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
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sdbeard/common-services/auth/secure"
	logger "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

// TODO: Replace with secret from secrets manager
//var secretKey = []byte("another-secret-key")

// Cookies
// https://golang.ch/how-to-work-with-cookies-in-golang/#:~:text=Basic%20usage%20of%20Cookies%20with%20Golang%201%20Name,SameSite%20constants%20from%20the%20net%2Fhttp%20package.%20More%20items

/**********************************************************************************/

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		render := render.New()

		authToken := getTokenFromSession(req)
		if authToken == "" {
			render.JSON(res, http.StatusUnauthorized, "no authorization information found")
			return
		}

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing")
			}
			secret, _ := secure.GetSecret("jwtsecretkey")
			return secret.Secret(), nil
		})
		if err != nil {
			render.JSON(res, http.StatusInternalServerError, err.Error())
			return
		}

		if !token.Valid {
			render.JSON(res, http.StatusUnauthorized, "not authorized")
			return
		}

		next.ServeHTTP(res, req)
	})
}

/**********************************************************************************/

func getTokenFromSession(req *http.Request) string {
	token, err := secure.GetSessionValue[string](req, "jwt")
	if err == nil {
		return token
	}

	// TODO:  Need to find a away to process the error
	return getTokenFromHeader(req)
}

func getTokenFromHeader(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return getTokenFromAuthCookie(req)
	}

	tokens := strings.Split(authHeader, "Bearer ")
	if len(tokens) == 2 {
		return tokens[1]
	}

	return ""
}

func getTokenFromAuthCookie(req *http.Request) string {
	cookie, err := req.Cookie("auth")
	if err != nil {
		logger.Error(err.Error())
		return ""
	}

	return cookie.Value
}

/**********************************************************************************/
