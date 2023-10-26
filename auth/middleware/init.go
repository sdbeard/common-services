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
package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sdbeard/common-services/auth/conf"
	"github.com/sdbeard/go-supportlib/common/util"
	"github.com/unrolled/render"
)

var (
	initfile = "auth.init"
)

/**********************************************************************************/

func IsInitialized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		render := render.New()

		if util.FileExists(fmt.Sprintf("%s%s%s", conf.Get().WorkingFolder, string(os.PathSeparator), initfile)) {
			next.ServeHTTP(res, req)
			return
		}

		if req.URL.Path != "/init" {
			render.JSON(res, http.StatusUnauthorized, "the service has not been initialized. please initialize to use the service")
			return
		}

		next.ServeHTTP(res, req)
	})
}

/**********************************************************************************/
