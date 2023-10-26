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
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sdbeard/common-services/auth/conf"
	"github.com/sdbeard/common-services/auth/middleware"
	"github.com/sdbeard/common-services/auth/secure"
	"github.com/sdbeard/common-services/auth/types"
	"github.com/sdbeard/go-supportlib/api/handlers"
	rest "github.com/sdbeard/go-supportlib/api/service"
	apitypes "github.com/sdbeard/go-supportlib/api/types"
	"github.com/sdbeard/go-supportlib/data/types/util/dataservice"
	logger "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

/**********************************************************************************/

func NewAuthService() (*AuthService, error) {
	newService := &AuthService{
		render: render.New(),
	}

	newService.RestService = rest.NewRestService(
		&conf.Get().ApiConf,
		newService.initializeRouter,
	)

	return newService, nil
}

/**********************************************************************************/

type AuthService struct {
	*rest.RestService
	render *render.Render
	//secret *sectypes.SimpleSecret
}

/***** exported functions *********************************************************/

// Start starts the running version of the API and is ready to receive requests
func (auth *AuthService) Start() error {
	defer func() {
		logger.Info("The hosting system has signaled the service to shutdown")
		auth.Stop()
	}()

	stopChannel := auth.createStopChannel()

	go auth.RestService.StartSimple()

	<-stopChannel
	close(stopChannel)

	return nil
}

// Stop initiaties the graceful shutdown of the API's underlying rest service
func (auth *AuthService) Stop() {
	auth.RestService.Stop()
}

/**********************************************************************************/

func (auth *AuthService) initializeRouter(router *mux.Router) {
	chain := alice.New(middleware.IsInitialized, handlers.LoggingHandler, handlers.JSONContentTypeHandler)
	//authChain := alice.New(middleware.IsInitialized, middleware.IsAuthorized, handlers.LoggingHandler, handlers.JSONContentTypeHandler)

	router.Handle("/", chain.Then(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		auth.render.JSON(res, http.StatusOK, "service called")
	})))

	router.Methods("GET").Path("/robots.txt").Handler(chain.ThenFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", mime.TypeByExtension(path.Ext("robots.txt")))
		res.WriteHeader(200)
		res.Write([]byte("# workshop-engine orgs service\n" +
			"# https://github.com/sdbeard/common-services/auth/\n" +
			"User-agent: *\n" +
			"Disallow: /"))
	}))

	apitypes.BaselineAPI(router, chain)

	router.Methods("POST").Path("/init").Handler(chain.ThenFunc(auth.enroll))
	//router.Methods("POST").Path("/authenticate").Handler(chain.ThenFunc(auth.authenticate))
	//router.Methods("GET").Path("/admin").Handler(authChain.ThenFunc(auth.adminIndex))
	//router.Methods("GET").Path("/user").Handler(authChain.ThenFunc(authapi.userIndex))
	//router.Methods("GET").Path("/index").Handler(alice.New().ThenFunc(authapi.index))

	auth.initializeSecretsRouter(router, chain)

	//router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8000")
	//	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, POST, PUT, OPTIONS")
	//	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header, Role")
	//})
}

func (auth *AuthService) createStopChannel() chan os.Signal {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGINT,
	)

	return stopChannel
}

func (auth *AuthService) enroll(res http.ResponseWriter, req *http.Request) {
	// Need to check if init file/flag has been set
	enrollment := new(types.Enrollment)

	err := json.NewDecoder(req.Body).Decode(&enrollment)
	if err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := secure.GenerateHashPassword(enrollment.User.Password)
	if err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}
	enrollment.User.Password = hashedPassword

	// Save the role, user and secret

	if err = dataservice.Add[*types.Role](dataservice.Request{
		Dataplane: conf.Get().Dataplane,
		Value:     enrollment.Role,
	}); err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	if err = dataservice.Add[*types.User](dataservice.Request{
		Dataplane: conf.Get().Dataplane,
		Value:     enrollment.User,
	}); err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	if err = auth.saveSecret(enrollment.JWTSecret); err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	if err = auth.saveSecret(enrollment.SessionSecret); err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	auth.render.JSON(res, http.StatusOK, "successfully initialized the service")
}

/*
func (auth *AuthService) authenticate(res http.ResponseWriter, req *http.Request) {
	authDetails := new(types.Authentication)

	err := json.NewDecoder(req.Body).Decode(authDetails)
	if err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	tags := dynamodb.ProcessDynamoTags(util.GetTypeObject[TUser]())
	hashKey, _ := tags.HashKey()
	if hashKey == "" {
		hashKey = conf.Get().Dataplane.Parameters["hashkey"].(string)
	}

	authUser, err := dataservice.GetItem[*types.User](dataservice.Request{
		Dataplane:  conf.Get().Dataplane,
		Key:        hashKey,
		Value:      authDetails.Username,
		Comparator: dsapi.EQ,
	})
	if err != nil {
		auth.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	check := secure.CheckPasswordHash(authDetails.Password, authUser.Password)
	if !check {
		auth.render.JSON(res, http.StatusUnauthorized, "username or password is incorrect")
		return
	}

	validToken, err := secure.GenerateJWT(auth.secret, authUser.Claims)
	if err != nil {
		auth.render.JSON(res, http.StatusUnauthorized, "failed to generate token")
		return
	}

	auth.render.JSON(res, http.StatusOK, validToken)
}

func (auth *AuthService) adminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}
	w.Write([]byte("Welcome, Admin."))
}
*/
/**********************************************************************************/
