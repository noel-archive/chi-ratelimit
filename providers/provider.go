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

// Package providers is the package to build your custom providers. We have
// built in support for etcd, Redis, and in-memory.
package providers

import "github.com/noelware/chi-ratelimit/types"

// Provider is the interface to implement the ratelimiting process itself.
type Provider interface {
	// Reset returns an error if there was an error of resetting
	// the types.Ratelimit object for the key provided.
	Reset(key string) (bool, error)

	// Close closes this Provider if needed.
	Close() error

	// Name returns the Provider name.
	Name() string

	// Put adds the key/value pair into the provider.
	Put(key string, value *types.Ratelimit) error

	// Get returns a tuple of the Ratelimit object/nil or an error if anything
	// has happened.
	//
	// The object can be nil if an error occurs, so you must add
	// appropriate error handling.
	Get(key string) (*types.Ratelimit, error)
}
