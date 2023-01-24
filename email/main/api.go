// **********************************************************************************
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
	"io/ioutil"
	"mime"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sdbeard/common-services/email/conf"
	"github.com/sdbeard/common-services/email/types"
	goapi "github.com/sdbeard/go-supportlib/api"
	"github.com/sdbeard/go-supportlib/api/handlers"
	rest "github.com/sdbeard/go-supportlib/api/service"
	apitypes "github.com/sdbeard/go-supportlib/api/types"
	"github.com/sdbeard/go-supportlib/common/util"
	"github.com/unrolled/render"
)

/**********************************************************************************/

// NewEmailAPI creates and returns a reference to a new EmailAPI struct
func NewEmailAPI() *EmailAPI {
	newAPI := &EmailAPI{
		render: render.New(),
	}

	newAPI.service = rest.NewRestService(
		&conf.GetConf().APIConf,
		newAPI.initializeRouter,
	)

	return newAPI
}

/***** EmailAPI *******************************************************************/

// EmailAPI is the representation of an email service API
type EmailAPI struct {
	router        *mux.Router
	service       *rest.RestService
	render        *render.Render
	worker        *types.EmailWorker
	isInitialized bool
}

/***** RESTService Interface Implementation ***************************************/

// GetState interrogates the service to understand its current state for status
// reporting purposes. This is a RESTService interface implementation
func (api *EmailAPI) GetState() goapi.State {
	return goapi.FUNCTIONAL
}

/***** exported functions *********************************************************/

// Start starts the running version of the API and is ready to receive requests
func (api *EmailAPI) Start() error {
	return api.service.StartSimple()
}

// Stop initiaties the graceful shutdown of the API's underlying rest service
func (api *EmailAPI) Stop() {
	api.service.Stop()
}

/**********************************************************************************/

func (api *EmailAPI) initializeRouter(router *mux.Router) {
	//stdChain := alice.New(handlers.LoggingHandler, handlers.JSONContentTypeHandler, auth.AuthMiddleware, sess.SessionMiddleware)
	chain := alice.New(handlers.LoggingHandler, handlers.JSONContentTypeHandler)

	router.Handle("/", chain.Then(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		api.render.JSON(res, http.StatusOK, "Service called")
	})))

	router.Methods("GET").Path("/robots.txt").Handler(chain.ThenFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", mime.TypeByExtension(path.Ext("robots.txt")))
		res.WriteHeader(200)
		res.Write([]byte("# workshop-engine orgs service\n" +
			"# https://github.com/sdbeard/workshop-engine/services/admin/\n" +
			"User-agent: *\n" +
			"Disallow: /"))
	}))

	emailRouter := router.PathPrefix("/kp/email").Subrouter()

	apitypes.BaselineAPI(emailRouter, chain)

	emailRouter.Methods("POST").Path("/{user}").Handler(chain.ThenFunc(api.sendEmail))
	emailRouter.Methods("GET").Path("/{user").Handler(chain.ThenFunc(api.getEmail))

	api.router = emailRouter
	api.isInitialized = true
}

/**********************************************************************************/

func (api EmailAPI) sendEmail(res http.ResponseWriter, req *http.Request) {
	// Get the user variable
	vars := mux.Vars(req)
	user, ok := vars["user"]
	if !ok {
		api.render.JSON(res, goapi.ErrMissingAttribute.HTTPCode, "the 'user' variable was not found")
	}
	_ = user

	// Get the email from the request body
	emailBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		api.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	// Get the email
	email, err := util.FromJSON[types.Email](emailBytes)
	if err != nil {
		api.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	if err = api.worker.SendEmail(email); err != nil {
		api.render.JSON(res, http.StatusInternalServerError, err.Error())
		return
	}

	api.render.JSON(res, http.StatusOK, "email successfully sent")
}

func (api EmailAPI) getEmail(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	user, ok := vars["user"]
	if !ok {
		api.render.JSON(res, goapi.ErrMissingAttribute.HTTPCode, "the 'user' variable was not found")
	}
	_ = user

	api.render.JSON(res, goapi.ErrNotImplemented.HTTPCode, "function not implemented")
}

/**********************************************************************************/
