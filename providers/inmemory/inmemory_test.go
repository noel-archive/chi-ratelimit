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

package inmemory

import (
	"github.com/noelware/chi-ratelimit/types"
	"testing"
	"time"
)

func TestProvider_Get(t *testing.T) {
	t.Log("Pre-init: initializing test TestProvider_Get")

	// have to add mock data
	provider := NewProvider()

	t.Log("Pre-init: Adding mock data...")
	err := provider.Put("owo", types.NewRatelimit(1000, false, time.Now().Add(1*time.Hour)))
	if err != nil {
		t.Errorf("unable to add object \"owo\" into provider: %s", err)
	}

	t.Log("Pre-init: Success! Now fetching...")

	// Now, let's check if it exists
	ratelimit, err := provider.Get("owo")
	if err != nil {
		t.Errorf("Unexpected error when retrieving \"owo\" object: %s", err)
	}

	t.Log("No errors had occurred while retrieving the \"owo\" object.")

	if ratelimit == nil {
		t.Error("Unexpected `nil` in ratelimit.")
	}

	t.Log("no unexpected nil values!")

	if ratelimit.Remaining != 999 {
		t.Errorf("Unexpected value %d (expected=999)", ratelimit.Remaining)
	}

	t.Log("Ratelimit#copy has been a success. :)")
}

func TestProvider_Reset(t *testing.T) {
	t.Log("Pre-init: initializing test TestProvider_Reset")

	// have to add mock data
	provider := NewProvider()

	t.Log("Pre-init: Adding mock data...")
	err := provider.Put("owo", types.NewRatelimit(1000, false, time.Now().Add(1*time.Hour)))
	if err != nil {
		t.Errorf("unable to add object \"owo\" into provider: %s", err)
	}

	t.Log("Pre-init: Success! Now resetting...")

	// Let's reset it
	reset, err := provider.Reset("owo")
	if err != nil {
		t.Errorf("Unexpected error in removing \"owo\" object: %s", err)
	}

	t.Log("No errors has occurred. :)")

	if !reset {
		t.Errorf("Unexpected value '%v' in removing \"owo\" object (expected=true)", reset)
	}

	t.Log("It worked as expected! :3")
}
