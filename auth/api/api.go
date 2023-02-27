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
package api

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sdbeard/common-services/auth/conf"
	"github.com/sdbeard/common-services/auth/middleware"
	"github.com/sdbeard/common-services/auth/secure"
	"github.com/sdbeard/common-services/auth/types"
	"github.com/sdbeard/go-supportlib/api/handlers"
	rest "github.com/sdbeard/go-supportlib/api/service"
	apitypes "github.com/sdbeard/go-supportlib/api/types"
	"github.com/sdbeard/go-supportlib/aws/service/dynamodb"
	"github.com/sdbeard/go-supportlib/common/util"
	"github.com/sdbeard/go-supportlib/data/types/dsapi"
	"github.com/sdbeard/go-supportlib/data/types/util/dataservice"
	"github.com/unrolled/render"
	//logger "github.com/sirupsen/logrus"
)

/**********************************************************************************/

func NewAuthApi[TUser types.AuthUser]() (*AuthApi[TUser], error) {
	newApi := &AuthApi[TUser]{
		render: render.New(),
		secret: []byte("secretkey"),
	}

	newApi.service = rest.NewRestService(
		&conf.Get().ApiConf,
		newApi.initializeRouter,
	)

	return newApi, nil
}

/**********************************************************************************/

type AuthApi[TUser types.AuthUser] struct {
	service *rest.RestService
	render  *render.Render
	secret  []byte
}

/***** exported functions *********************************************************/

// Start starts the running version of the API and is ready to receive requests
func (authapi *AuthApi[TUser]) Start() error {
	return authapi.service.StartSimple()
}

// Stop initiaties the graceful shutdown of the API's underlying rest service
func (authapi *AuthApi[TUser]) Stop() {
	authapi.service.Stop()
}

/**********************************************************************************/

func (authapi *AuthApi[TUser]) initializeRouter(router *mux.Router) {
	chain := alice.New(handlers.LoggingHandler, handlers.JSONContentTypeHandler)
	authChain := alice.New(handlers.LoggingHandler, middleware.IsAuthorized, handlers.JSONContentTypeHandler)

	router.Handle("/", chain.Then(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authapi.render.JSON(res, http.StatusOK, "service called")
	})))

	router.Methods("GET").Path("/robots.txt").Handler(chain.ThenFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", mime.TypeByExtension(path.Ext("robots.txt")))
		res.WriteHeader(200)
		res.Write([]byte("# workshop-engine orgs service\n" +
			"# https://github.com/sdbeard/workshop-engine/services/admin/\n" +
			"User-agent: *\n" +
			"Disallow: /"))
	}))

	apitypes.BaselineAPI(router, chain)

	router.Methods("POST").Path("/enroll").Handler(chain.ThenFunc(authapi.enroll))
	router.Methods("POST").Path("/authenticate").Handler(chain.ThenFunc(authapi.authenticate))
	router.Methods("GET").Path("/admin").Handler(authChain.ThenFunc(authapi.adminIndex))
	//router.Methods("GET").Path("/user").Handler(authChain.ThenFunc(authapi.userIndex))
	//router.Methods("GET").Path("/index").Handler(alice.New().ThenFunc(authapi.index))

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header, Role")
	})
}

func (authapi *AuthApi[TUser]) enroll(res http.ResponseWriter, req *http.Request) {
	user := util.GetTypeObject[TUser]()

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		authapi.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := secure.GenerateHashPassword(user.Password())
	if err != nil {
		authapi.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}
	user.SetPassword(hashedPassword)

	if err = dataservice.Add[TUser](dataservice.Request{
		Dataplane: conf.Get().Dataplane,
		Value:     user,
	}); err != nil {
		authapi.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	authapi.render.JSON(res, http.StatusOK, fmt.Sprintf("successfully added %s", user.UserId()))
}

func (authapi *AuthApi[TUser]) authenticate(res http.ResponseWriter, req *http.Request) {
	authDetails := new(types.Authentication)

	err := json.NewDecoder(req.Body).Decode(authDetails)
	if err != nil {
		authapi.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	tags := dynamodb.ProcessDynamoTags(util.GetTypeObject[TUser]())
	hashKey, _ := tags.HashKey()
	if hashKey == "" {
		hashKey = conf.Get().Dataplane.Parameters["hashkey"].(string)
	}

	authUser, err := dataservice.GetItem[TUser](dataservice.Request{
		Dataplane:  conf.Get().Dataplane,
		Key:        hashKey,
		Value:      authDetails.Username,
		Comparator: dsapi.EQ,
	})
	if err != nil {
		authapi.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	check := secure.CheckPasswordHash(authDetails.Password, authUser.Password())
	if !check {
		authapi.render.JSON(res, http.StatusUnauthorized, "username or password is incorrect")
		return
	}

	validToken, err := secure.GenerateJWT(authapi.secret, authUser.GetClaims())
	if err != nil {
		authapi.render.JSON(res, http.StatusUnauthorized, "failed to generate token")
		return
	}

	authapi.render.JSON(res, http.StatusOK, validToken)
}

func (authapi *AuthApi[TUser]) adminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}
	w.Write([]byte("Welcome, Admin."))
}

/**********************************************************************************/
