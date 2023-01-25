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
package types

import (
	"strings"

	"github.com/sdbeard/go-supportlib/aws"
	common "github.com/sdbeard/go-supportlib/common/types"
)

// <type>://<host>@<port>@<username>,<password>@<parameter map>
// smtp://smtp.gmail.com@587@kronedev@gmail.com,Password1
// ses://http://localstack.testlab.local@@Access Key,Secret Access Key@region=us-east-2

/**********************************************************************************/
/***** EmailConnectionConfig ******************************************************/

// EmailConnectionConfig represents the connection information required to connect
// to an email server/system
type EmailConnectionConfig struct {
	Credentials common.Credentials `json:"credentials,omitempty" dynamodbav:"credentials,omitmepty"`
	Parameters  map[string]string  `json:"parameters,omitempty" dynamodbav:"parameters,omitempty"`
	Host        string             `json:"host" dynamodbav:"host"`
	Port        string             `json:"port" dynamodbav:"port"`
}

/***** Marshal/Unmarshal interface implementations ********************************/

// UnmarshalText is the custom method use to create a EmailConnectionConfig when
// read from a string/byte slice or environment variable
func (config *EmailConnectionConfig) UnmarshalText(text []byte) error {
	newConnectConfig, err := EmailConnectionConfigFromString(string(text))
	if err != nil {
		return err
	}
	*config = newConnectConfig

	return nil
}

/***** exported functions *********************************************************/

// Creates and returns an AWS ConnectionConfig object ffrom the
// EmailConnectionConfig
func (config EmailConnectionConfig) AWSConnectConfig() aws.ConnectConfig {
	connectConfig := new(aws.ConnectConfig)
	connectConfig.Endpoint = config.Host
	connectConfig.Keys = aws.ConnectConfigKeys{
		AccessKeyID:     config.Credentials.Username,
		SecretAccessKey: config.Credentials.Password,
		Token:           config.Credentials.Token,
	}

	if region, ok := config.Parameters["region"]; ok {
		connectConfig.Region = region
	}

	return *connectConfig
}

/**********************************************************************************/

/***** exported functions *********************************************************/

// EmailConnectionConfigFromString parses the passed in connection string
// creates and returns a new EmailConnectionConfig object
func EmailConnectionConfigFromString(connection string) (EmailConnectionConfig, error) {
	var err error = nil
	newConfig := EmailConnectionConfig{}

	if connection != "" {
		err = getEmailHost(&newConfig, connection)
	}

	return newConfig, err
}

/**********************************************************************************/

func getEmailHost(emailConfig *EmailConnectionConfig, parameters string) error {
	parameter, parameterList := getNext(parameters)

	emailConfig.Host = parameter

	if parameterList == "" {
		return nil
	}

	return getEmailPort(emailConfig, parameterList)
}

func getEmailPort(emailConfig *EmailConnectionConfig, parameters string) error {
	parameter, parameterList := getNext(parameters)

	emailConfig.Port = parameter

	if parameterList == "" {
		return nil
	}

	return getCredentials(emailConfig, parameterList)
}

func getCredentials(emailConfig *EmailConnectionConfig, parameters string) error {
	parameter, parameterList := getNext(parameters)

	if parameter != "" {
		credentials := strings.Split(parameter, ",")
		getUsername(&emailConfig.Credentials, credentials)
	}

	if parameterList == "" {
		return nil
	}

	return getParameters(emailConfig, parameterList)
}
func getUsername(credentials *common.Credentials, credentialList []string) error {
	credentials.Username = credentialList[0]

	if len(credentialList) == 1 {
		return nil
	}

	return getPassword(credentials, credentialList[1:])
}

func getPassword(credentials *common.Credentials, credentialList []string) error {
	credentials.Password = credentialList[0]

	if len(credentialList) == 1 {
		return nil
	}

	return getAPIKey(credentials, credentialList[1:])
}

func getAPIKey(credentials *common.Credentials, credentialList []string) error {
	credentials.APIKey = credentialList[0]

	if len(credentialList) == 1 {
		return nil
	}

	return getToken(credentials, credentialList[1:])
}

func getToken(credentials *common.Credentials, credentialList []string) error {
	credentials.Token = credentialList[0]
	return nil
}

func getParameters(emailConfig *EmailConnectionConfig, parameters string) error {
	parameter, _ := getNext(parameters)
	connectionParameters := strings.Split(parameter, ",")

	if parameter == "" {
		return nil
	}

	emailConfig.Parameters = make(map[string]string)
	for _, connectionParameter := range connectionParameters {
		keyValue := strings.Split(connectionParameter, "=")
		if len(keyValue) != 2 {
			continue
		}
		emailConfig.Parameters[keyValue[0]] = keyValue[1]
	}

	return nil
}

func getNext(parameters string) (string, string) {
	if parameters == "" {
		return "", ""
	}

	first := strings.Index(parameters, "@")
	if first < 0 {
		// This is the last/only parameter
		return parameters, ""
	}

	second := strings.Index(parameters[first+1:], "@")
	if second < 0 {
		paramList := strings.Split(parameters, "@")
		return paramList[0], paramList[1]
	}

	if strings.Contains(parameters[:first+second+1], ",") {
		return parameters[:first+second+1], parameters[first+second+2:]
	}

	return parameters[:first], parameters[first+1:]

}

/**********************************************************************************/
