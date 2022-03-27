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

// Package inmemory implements the providers.Provider object for using
// an in-memory, mapped object. This isn't recommended for production,
// mainly for messing around and testing.
package inmemory

import (
	"fmt"
	"github.com/noelware/chi-ratelimit/providers"
	"github.com/noelware/chi-ratelimit/types"
)

type Provider struct {
	data map[string]*types.Ratelimit
}

func NewProvider() providers.Provider {
	return &Provider{
		data: make(map[string]*types.Ratelimit),
	}
}

func (p *Provider) Reset(key string) (bool, error) {
	// Check if the object exists
	var ratelimit *types.Ratelimit = nil

	for k, val := range p.data {
		if key == k {
			ratelimit = val
			break
		}
	}

	if ratelimit == nil {
		return false, nil
	}

	data := p.data
	delete(data, key)

	p.data = data
	return true, nil
}

func (*Provider) Close() error {
	return nil
}

func (*Provider) Name() string {
	return "in-memory provider"
}

func (p *Provider) Put(key string, value *types.Ratelimit) error {
	// Check if it exists
	for k, _ := range p.data {
		if key == k {
			return fmt.Errorf("ratelimit with key %s already exists", k)
		}
	}

	// Add it onto memory
	p.data[key] = value
	return nil
}

func (p *Provider) Get(key string) (*types.Ratelimit, error) {
	var ratelimit *types.Ratelimit
	for k, val := range p.data {
		if key == k {
			ratelimit = val
			break
		}
	}

	if ratelimit == nil {
		return nil, nil
	}

	// Update the inner map with the new ratelimit
	copied := ratelimit.Copy()
	data := p.data
	data[key] = copied

	p.data = data
	return copied, nil
}
