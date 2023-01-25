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
// **********************************************************************************
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/sdbeard/common-services/email/conf"
	"github.com/sdbeard/common-services/email/types"
	"github.com/sdbeard/go-supportlib/common/logging"
	logger "github.com/sirupsen/logrus"
)

var (
	configFile = flag.String("configfile", "config.json", "specifies the configuration file to use for the service configuration")
	build      = ""
	buildDate  = ""
	version    = "0.0.0"
)

/**********************************************************************************/

func init() {
	flag.Parse()

	// Load the configuration
	if err := conf.LoadConf(*configFile); err != nil {
		panic(err)
	}

	logging.InitializeLogging(conf.GetConf().LogConf)
}

func main() {
	fmt.Println("Workshop-Engine Orgs Service Init & Startup...")
	logger.WithFields(map[string]interface{}{
		"Version":    version,
		"Build":      build,
		"Build Date": buildDate,
		"GO Version": runtime.Version(),
		"PID":        os.Getpid(),
	}).Infof("Runtime configuration")

	config, err := types.EmailConnectionConfigFromString("")
	if err != nil {
		panic(err)
	}

	smtpClient := types.NewSmtpClient(config)
	if err = smtpClient.SendEmail(types.Email{}); err != nil {
		panic(err)
	}

	/*
		stopChannel := createStopChannel()

		api := NewEmailAPI()

		go api.Start()

		<-stopChannel
		close(stopChannel)

		// Stop the API
		api.Stop()
	*/
	logger.Info("Completed execution...shutting down")
}

/**********************************************************************************/

func createStopChannel() chan os.Signal {
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

/**********************************************************************************/
