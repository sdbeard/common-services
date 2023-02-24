// *********************************************************************************
// The MIT License (MIT)

// Copyright (c) 2023 Sean Beard

// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in the
// Software without restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the
// Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN
// AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// *********************************************************************************
package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sdbeard/env/v7"
	apicfg "github.com/sdbeard/go-supportlib/api/config"
	"github.com/sdbeard/go-supportlib/common/logging"
	"github.com/sdbeard/go-supportlib/common/util"
	"github.com/sdbeard/go-supportlib/data/types/configuration"
)

var config *Configuration

/***** Configuration **************************************************************/

// Configuration holds all of the necessary files for configuring an authentication
// and authorization service with JWTs
type Configuration struct {
	Dataplane     configuration.DataplaneConnection `json:"dataplane" env:"AUTH_DATAPLANE"`
	ApiConf       apicfg.ListenerConfig             `json:"api" env:"AUTH_APICONF"`
	LogConf       logging.LogConfig                 `json:"log" env:"AUTH_LOGCONF"`
	WorkingFolder string                            `json:"-"`
}

/***** exported functions *********************************************************/

// Get returns the reference to the current Configuration object
func Get() *Configuration {
	return config
}

// Load loads the configuration from thje value in the environment
func Load(file string) error {
	// Get the working folder
	workingDir, _ := os.Getwd()
	workingFolder, _ := filepath.Abs(workingDir)

	// Create the configuration object and set defaults
	config = new(Configuration)
	if file != "" && util.FileExists(file) {
		fileBytes, err := util.ReadFile(fmt.Sprintf("%s%s%s", workingFolder, string(os.PathSeparator), file))
		if err != nil {
			return err
		}

		// Unmarshal the configuration file
		err = json.Unmarshal(fileBytes, config)
		if err != nil {
			return err
		}
	}

	config.WorkingFolder = workingFolder

	// Load the environmnet variables from a file if they exist
	godotenv.Load()

	return env.Parse(config)
}

/**********************************************************************************/
