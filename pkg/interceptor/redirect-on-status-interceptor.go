// Package interceptor: contents of the file 'redirect_on_status_interceptor.go' were copied from:
// https://github.com/traefik/traefik/blob/master/pkg/middlewares/customerrors/custom_errors.go
// under the MIT license.

package interceptor

import (
	"bufio"
	"fmt"
	"github.com/m-riedel/traefik-plugin-redirect-on-status/pkg/types"
	"net"
	"net/http"
)

type RedirectOnStatusInterceptor struct {
	headerMap          http.Header
	responseWriter     http.ResponseWriter
	request            *http.Request
	code               int
	httpCodeRanges     types.HTTPCodeRanges
	caughtFilteredCode bool
	headersSent        bool
	redirectUri        string
	redirectCode       int
}

func NewRedirectOnStatusInterceptor(rw http.ResponseWriter, request *http.Request, redirectUri string, redirectCode int, httpCodeRanges types.HTTPCodeRanges) *RedirectOnStatusInterceptor {
	return &RedirectOnStatusInterceptor{
		headerMap:      make(http.Header),
		responseWriter: rw,
		code:           http.StatusOK,
		httpCodeRanges: httpCodeRanges,
		redirectUri:    redirectUri,
		request:        request,
		redirectCode:   redirectCode,
	}
}

func (r *RedirectOnStatusInterceptor) Header() http.Header {
	if r.headersSent {
		return r.responseWriter.Header()
	}

	if r.headerMap == nil {
		r.headerMap = make(http.Header)
	}

	return r.headerMap
}

func (r *RedirectOnStatusInterceptor) Write(buf []byte) (int, error) {
	// If WriteHeader was already called from the caller, this is a NOOP.
	// Otherwise, r.code is actually a 200 here.
	r.WriteHeader(r.code)

	if r.caughtFilteredCode {
		// We don't care about the contents of the response,
		// since we want to serve the ones from the error page,
		// so we just drop them.
		return len(buf), nil
	}
	return r.responseWriter.Write(buf)
}

func (r *RedirectOnStatusInterceptor) WriteHeader(code int) {
	if r.headersSent || r.caughtFilteredCode {
		return
	}

	// Handling informational headers.
	if code >= 100 && code <= 199 {
		// Multiple informational status codes can be used,
		// so here the copy is not appending the values to not repeat them.
		for k, v := range r.Header() {
			r.responseWriter.Header()[k] = v
		}

		r.responseWriter.WriteHeader(code)
		return
	}

	r.code = code

	if r.httpCodeRanges.Contains(r.code) {
		r.caughtFilteredCode = true
		r.code = r.redirectCode
		http.Redirect(r.responseWriter, r.request, r.redirectUri, r.code)
		return
	}

	// The copy is not appending the values,
	// to not repeat them in case any informational status code has been written.
	for k, v := range r.Header() {
		r.responseWriter.Header()[k] = v
	}

	r.responseWriter.WriteHeader(r.code)
	r.headersSent = true
}

// Hijack hijacks the connection.
func (r *RedirectOnStatusInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := r.responseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("%T is not a http.Hijacker", r.responseWriter)
}

// Flush sends any buffered data to the client.
func (r *RedirectOnStatusInterceptor) Flush() {
	// If WriteHeader was already called from the caller, this is a NOOP.
	// Otherwise, cc.code is actually a 200 here.
	r.WriteHeader(r.code)

	// We don't care about the contents of the response,
	// since we want to serve the ones from the error page,
	// so we just don't flush.
	// (e.g., To prevent superfluous WriteHeader on request with a
	// `Transfer-Encoding: chunked` header).
	if r.caughtFilteredCode {
		return
	}

	if flusher, ok := r.responseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
