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

import "fmt"

/**********************************************************************************/

// NewEmailWorker creates and returns a reference to a new EmailWorker instance
func NewEmailWorker() *EmailWorker {
	return &EmailWorker{}
}

/***** EmailWorker ****************************************************************/

// EmailWorker manages the operation of the email service and API
type EmailWorker struct{}

/***** exported functions *********************************************************/

// SendEmail takes the email as input and sends the email to the configured email
// source
func (worker EmailWorker) SendEmail(email Email) error {
	return fmt.Errorf("not implemented exception")
}

/**********************************************************************************/
/**********************************************************************************/
