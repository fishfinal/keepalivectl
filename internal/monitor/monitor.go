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

package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/fishfinal/keepalivectl/internal/config"
	"github.com/logrusorgru/aurora"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
)

// Monitor is responsible for establishing gRPC connections and monitoring their states.
type Monitor struct {
	aurora aurora.Aurora
	cfg    *config.Config
}

// Result represents a single monitoring result from a connection.
type Result struct {
	ConnectionID int
	Timestamp    time.Time
	State        connectivity.State
	StateChanged bool
	Message      string
}

// NewMonitor creates a new Monitor instance with the given configuration.
func NewMonitor(cfg *config.Config) *Monitor {
	return &Monitor{
		cfg: cfg,
	}
}

// StartConnection establishes a gRPC connection and monitors its state changes.
// It sends results to the provided channel until the context is cancelled.
func (m *Monitor) StartConnection(ctx context.Context, id int, resultCh chan<- Result) {
	kacp := keepalive.ClientParameters{
		Time:                m.cfg.KeepaliveTime,
		Timeout:             m.cfg.KeepaliveTimeout,
		PermitWithoutStream: true,
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
	}

	conn, err := grpc.DialContext(ctx, m.cfg.Endpoint, opts...)
	if err != nil {
		resultCh <- Result{
			ConnectionID: id,
			Timestamp:    time.Now(),
			Message:      fmt.Sprintf("connection failed: %v", err),
		}
		return
	}
	defer conn.Close()

	resultCh <- Result{
		ConnectionID: id,
		Timestamp:    time.Now(),
		Message: fmt.Sprintf("connected successfully, monitoring (keepalive: %v, timeout: %v)",
			m.cfg.KeepaliveTime, m.cfg.KeepaliveTimeout),
	}

	ticker := time.NewTicker(m.cfg.CheckInterval)
	defer ticker.Stop()

	var lastState connectivity.State
	stateChangeCount := 0

	for {
		select {
		case <-ctx.Done():
			resultCh <- Result{
				ConnectionID: id,
				Timestamp:    time.Now(),
				Message:      fmt.Sprintf("stopped monitoring, total state changes: %d", stateChangeCount),
			}
			return
		case <-ticker.C:
			state := conn.GetState()
			stateChanged := state != lastState
			if stateChanged {
				stateChangeCount++
				lastState = state
			}

			// Only output detailed information when state changes to keep output clean
			msg := fmt.Sprintf("state: %s", state)
			if stateChanged {
				msg = fmt.Sprintf("state changed: %s (change #%d)", state, stateChangeCount)
			}

			resultCh <- Result{
				ConnectionID: id,
				Timestamp:    time.Now(),
				State:        state,
				StateChanged: stateChanged,
				Message:      msg,
			}

			if state == connectivity.TransientFailure {
				resultCh <- Result{
					ConnectionID: id,
					Timestamp:    time.Now(),
					Message:      "⚠️ connection entered TransientFailure state, keepalive may be failing!",
				}
			}
		}
	}
}
