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
	"time"

	"github.com/golang-jwt/jwt/v5"
	logger "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

// TODO: Replace with secret from secrets manager
var secretKey = []byte("your-secret-key")

// Cookies
// https://golang.ch/how-to-work-with-cookies-in-golang/#:~:text=Basic%20usage%20of%20Cookies%20with%20Golang%201%20Name,SameSite%20constants%20from%20the%20net%2Fhttp%20package.%20More%20items

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		render := render.New()

		authToken := getTokenFromHeader(req)
		if authToken == "" {
			render.JSON(res, http.StatusUnauthorized, "no authorization information found")
			return
		}

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing")
			}
			return secretKey, nil
		})
		if err != nil {
			render.JSON(res, http.StatusInternalServerError, "your token has expired")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			req.Header.Set("Role", claims["role"].(string))
			next.ServeHTTP(res, req)
			return
		}

		render.JSON(res, http.StatusUnauthorized, "not authorized")
	})
}

/**********************************************************************************/

func generateJWT(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"roles": roles,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
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
