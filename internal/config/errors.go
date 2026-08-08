// Copyright (c) 2026 fishfinal
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package config

import "errors"

// Error definitions for configuration validation.
var (
	// ErrEmptyEndpoint indicates that the endpoint address is empty.
	ErrEmptyEndpoint = errors.New("endpoint cannot be empty")

	// ErrInvalidKeepaliveTime indicates that the keepalive time is invalid.
	ErrInvalidKeepaliveTime = errors.New("keepalive time must be greater than 0")

	// ErrInvalidKeepaliveTimeout indicates that the keepalive timeout is invalid.
	ErrInvalidKeepaliveTimeout = errors.New("keepalive timeout must be greater than 0")

	// ErrInvalidDuration indicates that the test duration is invalid.
	ErrInvalidDuration = errors.New("duration must be greater than 0")

	// ErrInvalidCheckInterval indicates that the check interval is invalid.
	ErrInvalidCheckInterval = errors.New("check interval must be greater than 0")

	// ErrInvalidConcurrency indicates that the concurrency count is invalid.
	ErrInvalidConcurrency = errors.New("concurrency must be greater than 0")
)
