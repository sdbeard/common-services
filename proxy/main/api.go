// *********************************************************************************
// The MIT License (MIT)
//
// # Copyright (c) 2024 Sean Beard
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
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sdbeard/common-services/proxy/conf"
	"github.com/sdbeard/go-supportlib/api/handlers"
	rest "github.com/sdbeard/go-supportlib/api/service"
	apitypes "github.com/sdbeard/go-supportlib/api/types"
	"github.com/sdbeard/go-supportlib/common/logging"
	logger "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

/**********************************************************************************/
// NewSCMWebhook create and returns a reference to a new SCMWebhook object
func NewProxy() *Proxy {
	logger.WithFields(logging.LogEntryContext(logger.Fields{})).Debug("")

	newProxy := &Proxy{
		render: render.New(),
	}

	newProxy.RestService = rest.NewRestService(
		&conf.Get().APIConfig,
		newProxy.initializeRouter,
	)

	return newProxy
}

/***** Proxy **********************************************************************/

// Proxy defines the api endpoints to receive requests to be proxied
type Proxy struct {
	*rest.RestService
	render *render.Render
}

/***** exported functions *********************************************************/

// Start starts the running version of the API and makes ready to receive requests
func (proxy *Proxy) Start() error {
	logger.WithFields(logging.LogEntryContext(logger.Fields{})).Debug("")

	defer func() {
		logger.Info("The hosting system has signaled the service to shutdown")
		proxy.Stop()
	}()

	stop := proxy.createStopChannel()

	go proxy.RestService.StartSimple()

	<-stop
	close(stop)

	return nil
}

// Stop initiaties the graceful shutdown of the API's underlying rest service
func (proxy *Proxy) Stop() {
	logger.WithFields(logging.LogEntryContext(logger.Fields{})).Debug("")

	proxy.RestService.Stop()
}

/**********************************************************************************/

func (proxy *Proxy) initializeRouter(router *mux.Router) {
	logger.WithFields(logging.LogEntryContext(logger.Fields{})).Debug("")

	chain := alice.New(handlers.LoggingHandler, handlers.JSONContentTypeHandler, handlers.GzipHandler)

	router.Methods("GET").Path("/robots.txt").Handler(chain.ThenFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", mime.TypeByExtension(path.Ext("robots.txt")))
		res.WriteHeader(200)
		res.Write([]byte("# workshop-engine orgs service\n" +
			"# https://github.com/sdbeard/common-services/auth/\n" +
			"User-agent: *\n" +
			"Disallow: /"))
	}))

	apitypes.BaselineAPI(router, chain)
}

func (proxy *Proxy) createStopChannel() chan os.Signal {
	logger.WithFields(logging.LogEntryContext(logger.Fields{})).Debug("")

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

/***** request handlers functions *************************************************/
/**********************************************************************************/
