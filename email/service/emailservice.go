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
package service

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sdbeard/common-services/email/api"
	logger "github.com/sirupsen/logrus"
)

/**********************************************************************************/

// NewEmailService creates and returns a reference to a new EmailService instance
func NewEmailService() *EmailService {
	return &EmailService{}
}

/***** EmailService ***************************************************************/

// EmailService manages the operation of the email service and API
type EmailService struct {
	api           *api.EmailAPI
	stopChannel   chan os.Signal
	isInitialized bool
}

/***** exported functions *********************************************************/

// Start initializes and starts the Email service
func (svc *EmailService) Start() error {
	// Set the defer function to close out the service
	defer func() {
		logger.Info("The hosting system has signaled the service to shutdown")
		svc.api.Stop()
	}()

	// Initialize the service
	svc.initialize()

	// Start the Admin services API
	go svc.api.Start()

	<-svc.stopChannel
	close(svc.stopChannel)

	return nil
}

/**********************************************************************************/

func (svc *EmailService) initialize() {
	defer func() { svc.isInitialized = true }()

	if svc.isInitialized {
		return
	}

	// Create the top-level API
	svc.api = api.NewEmailAPI()

	// Create the stop channel
	svc.createStopChannel()
}

func (svc *EmailService) createStopChannel() {
	svc.stopChannel = make(chan os.Signal, 1)
	signal.Notify(svc.stopChannel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGINT,
	)
}

/**********************************************************************************/
