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

package types

import "time"

// Ratelimit is the main struct that contains metadata about it.
type Ratelimit struct {
	// ResetTime returns a time.Time instance of when this Ratelimit
	// is reset and should be cleared from cache.
	ResetTime time.Time `json:"reset_time"`

	// Remaining returns an int32 of how many requests are left.
	Remaining int32 `json:"remaining"`

	// Global returns a bool if this is from a global route rather a scoped
	// one.
	Global bool `json:"global"`

	// Limit returns an int32 of the limit of requests this Ratelimit has.
	Limit int32 `json:"limit"`
}

// NewRatelimit creates a new Ratelimit pointer object.
func NewRatelimit(limit int32, global bool, resetTime time.Time) *Ratelimit {
	return &Ratelimit{
		ResetTime: resetTime,
		Remaining: limit,
		Global:    global,
		Limit:     limit,
	}
}

// Copy copies the current Ratelimit object into a new Ratelimit object
// but removing the remaining counter by one.
func (r *Ratelimit) Copy() *Ratelimit {
	this := r

	if r.Remaining > 0 {
		this.Remaining = r.Remaining - 1
	} else {
		this.Remaining = 0
	}

	return this
}

// Expired returns a bool if this Ratelimit has been expired or not.
func (r *Ratelimit) Expired() bool {
	return r.ResetTime.UnixNano() >= time.Now().UnixNano()
}

// Exceeded returns a bool if this Ratelimit has exceeded the rate limit.
func (r *Ratelimit) Exceeded() bool {
	return !r.Expired() && r.Remaining == 0
}
