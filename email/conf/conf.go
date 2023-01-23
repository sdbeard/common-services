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
package conf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdbeard/env"
	apicfg "github.com/sdbeard/go-supportlib/api/config"
	"github.com/sdbeard/go-supportlib/common/logging"
	"github.com/sdbeard/go-supportlib/common/util"
)

var (
	config *Configuration
)

/***** Configuration **************************************************************/

// Configuration holds the relevant values to configure the overall UtilData service
type Configuration struct {
	APIConf          apicfg.ListenerConfig `json:"api" env:"EMAIL_APICONF"`
	ConnectionString interface{}           `json:"connection" env:"EMAIL_CONNECTION"`
	LogConf          logging.LogConfig     `json:"log" env:"EMAIL_LOGCONF"`
	WorkingFolder    string                `json:"-"`
}

/**********************************************************************************/

/***** exported functions *********************************************************/

// Get returns the reference to the current configuration object
func GetConf() *Configuration {
	return config
}

func LoadConf(file string) error {
	// Get the working folder
	workingDir, _ := os.Getwd()
	workingFolder, _ := filepath.Abs(workingDir)

	// Create the configuration object and set defaults
	config = &Configuration{}
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

	// Load the environ file it it exists
	loadEnvironFile(workingFolder)

	return env.Parse(config)
}

/**********************************************************************************/

// loadEnvironFile looks for an optional file in the current execution folder and
// loads the contents of the '.environ' file into the local execution environment.
func loadEnvironFile(workingFolder string) error {
	// Retrieve and load the .environ file if it exists in the current folder
	if _, err := os.Stat(fmt.Sprintf("%s%s.environ", workingFolder, string(os.PathSeparator))); err != nil {
		// File doesn;t exists dimply return
		return nil
	}

	environFile, err := os.Open(fmt.Sprintf("%s%s.environ",
		workingFolder, string(os.PathSeparator)))
	if err != nil {
		return err
	}
	defer environFile.Close()

	// Load the environment from the local execution environment
	scanner := bufio.NewScanner(environFile)

	for scanner.Scan() {
		// Get the environment variable
		envVar := strings.Split(scanner.Text(), "=")
		if len(envVar) == 2 {
			err = os.Setenv(envVar[0], envVar[1])
			if err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

/**********************************************************************************/
