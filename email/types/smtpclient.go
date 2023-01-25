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
	"fmt"
	"net/smtp"
)

// https://www.geeksforgeeks.org/sending-email-using-smtp-in-golang/

func NewSmtpClient(config EmailConnectionConfig) *SmtpClient {
	return &SmtpClient{
		config: config,
	}
}

type SmtpClient struct {
	config EmailConnectionConfig
}

func (client SmtpClient) SendEmail(email Email) error {
	return smtp.SendMail(
		fmt.Sprintf("%s:%s", client.config.Host, client.config.Port),
		smtp.PlainAuth(
			"",
			client.config.Credentials.Username,
			client.config.Credentials.Password,
			client.config.Host,
		),
		"kronedev@gmail.com",
		[]string{"sean.beard@parivedasolutions.com"},
		[]byte("Hello from golang!"),
	)
}
