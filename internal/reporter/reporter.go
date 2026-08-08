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

package reporter

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/shaichunfeng/gologger"

	"github.com/fishfinal/keepalivectl/internal/config"
	"github.com/fishfinal/keepalivectl/internal/monitor"
	"google.golang.org/grpc/connectivity"
)

// Reporter collects and reports monitoring results.
type Reporter struct {
	totalChecks  int64
	readyCount   int64
	failureCount int64
	stateChanges int64
	colorful     aurora.Aurora
}

// NewReporter creates a new Reporter instance with color output enabled.
func NewReporter() *Reporter {
	return &Reporter{
		colorful: aurora.NewAurora(true),
	}
}

// Report processes a single monitoring result and outputs it to the log.
func (r *Reporter) Report(result monitor.Result) {
	atomic.AddInt64(&r.totalChecks, 1)

	if result.StateChanged {
		atomic.AddInt64(&r.stateChanges, 1)
	}

	if result.State == connectivity.Ready {
		atomic.AddInt64(&r.readyCount, 1)
	}

	if result.State == connectivity.TransientFailure {
		atomic.AddInt64(&r.failureCount, 1)
	}

	// Format output
	prefix := fmt.Sprintf("[Connection %d]", result.ConnectionID)
	stateStr := getStateString(result.State, r.colorful)

	// Output using gologger
	var msg string
	switch {
	case result.StateChanged:
		msg = fmt.Sprintf("%s %s %s %s",
			r.colorful.Cyan(result.Timestamp.Format(time.RFC3339)),
			r.colorful.Yellow(prefix),
			stateStr,
			r.colorful.BrightBlue(result.Message),
		)
	case result.State == connectivity.TransientFailure:
		msg = fmt.Sprintf("%s %s %s %s",
			r.colorful.Cyan(result.Timestamp.Format(time.RFC3339)),
			r.colorful.Yellow(prefix),
			stateStr,
			r.colorful.Red(result.Message),
		)
	default:
		msg = fmt.Sprintf("%s %s %s %s",
			r.colorful.Cyan(result.Timestamp.Format(time.RFC3339)),
			r.colorful.Yellow(prefix),
			stateStr,
			result.Message,
		)
	}

	gologger.Info().Msg(msg)
}

// PrintSummary prints a summary of all test results.
func (r *Reporter) PrintSummary(cfg *config.Config) {
	separator := strings.Repeat("=", 60)
	cyan := r.colorful.BrightCyan
	blue := r.colorful.BrightBlue

	gologger.Info().Msg("")
	gologger.Info().Msg(cyan(separator).String())
	gologger.Info().Msg(cyan("📊 Test Summary").String())
	gologger.Info().Msg(cyan(separator).String())
	gologger.Info().Msgf("%s %s", blue("Endpoint:"), cfg.Endpoint)
	gologger.Info().Msgf("%s %v", blue("Keepalive Interval:"), cfg.KeepaliveTime)
	gologger.Info().Msgf("%s %v", blue("Keepalive Timeout:"), cfg.KeepaliveTimeout)
	gologger.Info().Msgf("%s %v", blue("Test Duration:"), cfg.Duration)
	gologger.Info().Msgf("%s %d", blue("Concurrency:"), cfg.Concurrency)
	gologger.Info().Msg(cyan(separator).String())
	gologger.Info().Msgf("%s %d", blue("Total Checks:"), atomic.LoadInt64(&r.totalChecks))
	gologger.Info().Msgf("%s %d", blue("Ready State Count:"), atomic.LoadInt64(&r.readyCount))
	gologger.Info().Msgf("%s %d", blue("State Changes:"), atomic.LoadInt64(&r.stateChanges))

	failures := atomic.LoadInt64(&r.failureCount)
	if failures > 0 {
		gologger.Error().Msgf("❌ Failures: %d", failures)
	} else {
		gologger.Info().Msg("✅ All connections healthy, server Keepalive support is good")
	}
	gologger.Info().Msg(cyan(separator).String())
}

// getStateString returns a colored string representation of the connection state.
func getStateString(state connectivity.State, au aurora.Aurora) string {
	switch state {
	case connectivity.Ready:
		return au.Green("READY").String()
	case connectivity.Connecting:
		return au.Yellow("CONNECTING").String()
	case connectivity.TransientFailure:
		return au.Red("TRANSIENT_FAILURE").String()
	case connectivity.Idle:
		return au.Blue("IDLE").String()
	case connectivity.Shutdown:
		return au.Magenta("SHUTDOWN").String()
	default:
		return state.String()
	}
}
