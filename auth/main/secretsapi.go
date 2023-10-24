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
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sdbeard/common-services/auth/conf"
	"github.com/sdbeard/common-services/auth/types"
	"github.com/sdbeard/go-supportlib/secure/secrets/factory"
)

/**********************************************************************************/

func (auth *AuthService) initializeSecretsRouter(router *mux.Router, chain alice.Chain) {
	secretsRouter := router.PathPrefix("/secrets").Subrouter()

	//secretsRouter.Methods("GET").Path("").Handler(chain.ThenFunc(auth.getSecrets))
	//router.Methods("GET").Path("/secrets").Handler(authChain.ThenFunc(auth.getSecrets))
	//router.Methods("POST").Path("/secrets").Handler(authChain.ThenFunc(auth.addSecret))
	secretsRouter.Methods("GET").Path("/{secretid}").Handler(chain.ThenFunc(auth.getSecret))
	//router.Methods("PUT").Path("/secrets/{secretid}").Handler(chain.ThenFunc(auth.updateSecret))
	//router.Methods("GET").Path("/user").Handler(authChain.ThenFunc(authapi.userIndex))
	//router.Methods("GET").Path("/index").Handler(alice.New().ThenFunc(authapi.index))
}

func (auth *AuthService) getSecret(res http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.RemoteAddr, "localhost") && strings.Contains("", "localhost") {
		//Allow CORS here By * or specific origin
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	vars := mux.Vars(req)
	secretId := vars["secretid"]

	manager, err := factory.SecretsManagerFactory[*types.JWTSecret](conf.Get().SecretsConf)
	if err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	secret, err := manager.Retrieve(
		manager.Retrieve.WithSecretName(secretId),
	)
	if err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	auth.render.JSON(res, http.StatusOK, secret)
}

/**********************************************************************************/
