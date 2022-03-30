// 🔬 chi-ratelimit: Simple production-ready ratelimiter for go-chi applications
// Copyright (c) 2022 Noelware
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package ratelimit is the package to create production-ready ratelimiting
// in go-chi applications; more broad would be "net/http" servers.
package ratelimit

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/noelware/chi-ratelimit/providers"
	"github.com/noelware/chi-ratelimit/providers/inmemory"
	"github.com/noelware/chi-ratelimit/types"
)

// ~ ratelimiter options ~ \\

type Options struct {
	// CheckIfGlobalFunc returns a function if the given ratelimit is a global one.
	CheckIfGlobalFunc func(w http.ResponseWriter, req *http.Request) bool

	// OnRatelimit returns a function when the ratelimit has been reached. This is called
	// after the headers of `Retry-After` was applied, so you don't need to. :)
	OnRatelimit func(w http.ResponseWriter, req *http.Request)

	// KeyFunc returns a function for the ratelimit key that is stored. By default,
	// it will use the IP given to us.
	KeyFunc func(w http.ResponseWriter, req *http.Request) string

	// DefaultTimeWindow is a time.Duration instance of the default window
	// a ratelimit should be exceeded for.
	DefaultTimeWindow time.Duration

	// DefaultLimit returns an int of the default limit before it halts
	// if it has exceeded.
	DefaultLimit int

	// Provider is a providers.Provider instance to provide a persistent
	// data-layer in this Ratelimiter.
	Provider providers.Provider
}

type OptionsOverrideFunc func(o *Options)

// WithCheckIfGlobalFunc overrides the default option of the CheckIfGlobal option.
func WithCheckIfGlobalFunc(fun func(w http.ResponseWriter, req *http.Request) bool) OptionsOverrideFunc {
	return func(o *Options) {
		o.CheckIfGlobalFunc = fun
	}
}

// WithDefaultTimeWindow overrides the default option of the DefaultTimeWindow option.
func WithDefaultTimeWindow(t time.Duration) OptionsOverrideFunc {
	return func(o *Options) {
		o.DefaultTimeWindow = t
	}
}

// WithDefaultLimit overrides the default option of the DefaultLimit option.
func WithDefaultLimit(limit int) OptionsOverrideFunc {
	return func(o *Options) {
		o.DefaultLimit = limit
	}
}

// WithProvider overrides the default option of the Provider option.
func WithProvider(provider providers.Provider) OptionsOverrideFunc {
	return func(o *Options) {
		o.Provider = provider
	}
}

// WithKeyFunc overrides the default option of the `KeyFunc` option.
func WithKeyFunc(fun func(w http.ResponseWriter, req *http.Request) string) OptionsOverrideFunc {
	return func(o *Options) {
		o.KeyFunc = fun
	}
}

// WithOnRatelimit overrides the default option of the `OnRatelimit` option.
func WithOnRatelimit(fun func(w http.ResponseWriter, req *http.Request)) OptionsOverrideFunc {
	return func(o *Options) {
		o.OnRatelimit = fun
	}
}

func defaultOptions() *Options {
	return &Options{
		DefaultTimeWindow: 1 * time.Hour,
		DefaultLimit:      100,
		Provider:          inmemory.NewProvider(),
		CheckIfGlobalFunc: func(w http.ResponseWriter, req *http.Request) bool {
			return true
		},

		KeyFunc: func(w http.ResponseWriter, req *http.Request) string {
			return realIP(req)
		},

		OnRatelimit: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		},
	}
}

// Ratelimiter is the main struct that you must use to run a production-ready
// ratelimiting system into net/http servers.
type Ratelimiter struct {
	// checkIfGlobalFunc returns a function if the given ratelimit is a global one.
	checkIfGlobalFunc func(w http.ResponseWriter, req *http.Request) bool

	// onRatelimit returns a function when the ratelimit has been reached. This is called
	// after the headers of `Retry-After` was applied, so you don't need to. :)
	onRatelimit func(w http.ResponseWriter, req *http.Request)

	// keyFunc returns a function for the ratelimit key that is stored. By default,
	// it will use the IP given to us.
	keyFunc func(w http.ResponseWriter, req *http.Request) string

	// defaultTimeWindow is a time.Duration instance of the default window
	// a ratelimit should be exceeded for.
	defaultTimeWindow time.Duration

	// defaultLimit returns an int of the default limit before it halts
	// if it has exceeded.
	defaultLimit int

	// provider is a providers.Provider instance to provide a persistent
	// data-layer in this Ratelimiter.
	provider providers.Provider
}

func NewRatelimiter(opts ...OptionsOverrideFunc) *Ratelimiter {
	options := defaultOptions()
	for _, override := range opts {
		override(options)
	}

	return &Ratelimiter{
		checkIfGlobalFunc: options.CheckIfGlobalFunc,
		defaultTimeWindow: options.DefaultTimeWindow,
		defaultLimit:      options.DefaultLimit,
		keyFunc:           options.KeyFunc,
		onRatelimit:       options.OnRatelimit,
		provider:          options.Provider,
	}
}

func (r *Ratelimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Define global constants
		isGlobal := r.checkIfGlobalFunc(w, req)
		key := r.keyFunc(w, req)
		rl, err := r.provider.Get(key)
		headers := w.Header()

		if rl != nil {
			rl = types.NewRatelimit(
				int32(r.defaultLimit),
				isGlobal,
				time.Now().Add(r.defaultTimeWindow),
			)

			r.provider.Put(key, rl)
		}

		// TODO: add a option to call on error handling?
		if err != nil {
			panic(err)
		}

		if rl.Exceeded() {
			retry := strconv.FormatInt(time.Now().Sub(rl.ResetTime).Milliseconds(), 10)
			headers.Set("Retry-After", retry)

			r.onRatelimit(w, req)
			return
		}

		isGlobalString := ""
		if isGlobal {
			isGlobalString = "true"
		} else {
			isGlobalString = "false"
		}

		headers.Set("X-RateLimit-Remaining", strconv.Itoa(int(rl.Remaining)))
		headers.Set("X-RateLimit-Global", isGlobalString)
		headers.Set("X-RateLimit-Reset", strconv.FormatInt(rl.ResetTime.Unix()*1000, 10))
		headers.Set("X-RateLimit-Limit", strconv.Itoa(int(rl.Limit)))

		next.ServeHTTP(w, req)
	})
}

// https://github.com/go-chi/httprate/blob/master/httprate.go#L25-L47

func realIP(req *http.Request) string {
	var ip string
	if tcip := req.Header.Get("True-Client-IP"); tcip != "" {
		ip = tcip
	} else if xrip := req.Header.Get("X-Real-IP"); xrip != "" {
		ip = xrip
	} else if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		idx := strings.Index(xff, ", ")
		if idx == -1 {
			idx = len(xff)
		}

		// python moment
		ip = xff[:idx]
	} else {
		var err error

		ip, _, err = net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			ip = req.RemoteAddr
		}
	}

	return ip
}
