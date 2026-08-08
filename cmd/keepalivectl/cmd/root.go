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

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fishfinal/keepalivectl/cmd/keepalivectl/helper/flagshelper"
	"github.com/fishfinal/keepalivectl/internal/config"
	"github.com/fishfinal/keepalivectl/internal/monitor"
	"github.com/fishfinal/keepalivectl/internal/reporter"
	"github.com/spf13/cobra"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "keepalivectl",
	Short: "gRPC Keepalive connection testing tool",
	Long: `keepalivectl is a command-line tool for testing gRPC server Keepalive and connection pooling features.

It validates server-side keepalive handling by establishing long-lived connections and sending Keepalive pings,
with support for concurrent connection testing, real-time connection state monitoring, and detailed statistics reporting.

Core Features:
  • Establish gRPC long-lived connections and send Keepalive pings
  • Monitor connection state changes (Idle, Ready, TransientFailure)
  • Support concurrent connection testing
  • Customizable Keepalive parameters (interval, timeout)
  • Real-time reporting and final statistics summary`,
	Example: `  # Basic test - connect to the specified endpoint, run for 2 minutes
  keepalivectl -e localhost:9600 -d 2m

  # Concurrent test - 10 concurrent connections, run for 3 minutes
  keepalivectl -c 10 -d 3m

  # Custom Keepalive parameters - 30s interval, 10s timeout
  keepalivectl -t 30s -T 10s -d 5m

  # Quick test - 30 seconds run, 1 second check interval
  keepalivectl -d 30s -i 1s

  # High concurrency test - 50 connections, 5 minutes run
  keepalivectl -c 50 -d 5m -i 5s

  # Custom server address with full parameters
  keepalivectl -e grpc-server:9600 -c 5 -d 2m -t 15s -T 5s -i 2s`,
	Version: "1.0.0",
	RunE:    runKeepaliveTest,
}

const (
	targetGroup            = "Target"
	testBehaviorGroup      = "Test Behavior"
	keepaliveSettingsGroup = "Keepalive Settings"
)

func init() {
	cfg = config.New()

	// Disable automatic sorting to preserve the order of flags
	rootCmd.Flags().SortFlags = false

	// ===== Target Group =====
	rootCmd.Flags().StringVarP(&cfg.Endpoint, "endpoint", "e", "localhost:9600", "gRPC server endpoint address")
	rootCmd.Flags().SetAnnotation("endpoint", "group", []string{targetGroup})

	// ===== Test Behavior Group =====
	rootCmd.Flags().DurationVarP(&cfg.Duration, "duration", "d", 2*time.Minute, "Test duration")
	rootCmd.Flags().SetAnnotation("duration", "group", []string{testBehaviorGroup})

	rootCmd.Flags().DurationVarP(&cfg.CheckInterval, "check-interval", "i", 3*time.Second, "Connection state check interval")
	rootCmd.Flags().SetAnnotation("check-interval", "group", []string{testBehaviorGroup})

	rootCmd.Flags().IntVarP(&cfg.Concurrency, "concurrency", "c", 1, "Number of concurrent connections")
	rootCmd.Flags().SetAnnotation("concurrency", "group", []string{testBehaviorGroup})

	// ===== Keepalive Settings Group =====
	rootCmd.Flags().DurationVarP(&cfg.KeepaliveTime, "keepalive-time", "t", 10*time.Second, "Keepalive ping interval")
	rootCmd.Flags().SetAnnotation("keepalive-time", "group", []string{keepaliveSettingsGroup})

	rootCmd.Flags().DurationVarP(&cfg.KeepaliveTimeout, "keepalive-timeout", "T", 5*time.Second, "Keepalive ping timeout")
	rootCmd.Flags().SetAnnotation("keepalive-timeout", "group", []string{keepaliveSettingsGroup})

	// ===== Custom help template =====
	flagHelper := flagshelper.NewFlagsHelper(flagshelper.WithSpacing(24), flagshelper.WithShortFlag(true))
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", cmd.Long)
		fmt.Fprintf(cmd.OutOrStdout(), "Usage:\n  %s [flags]\n\n", cmd.CommandPath())

		flagHelper.PrintFlagGroup(cmd, targetGroup)
		flagHelper.PrintFlagGroup(cmd, testBehaviorGroup)
		flagHelper.PrintFlagGroup(cmd, keepaliveSettingsGroup)

		fmt.Fprintf(cmd.OutOrStdout(), "\nGlobal Flags:\n  -h, --help   help for %s\n", cmd.Name())
	})
}

// Execute runs the root command and returns any error that occurs.
func Execute() error {
	return rootCmd.Execute()
}

// runKeepaliveTest is the main entry point for the keepalive test.
func runKeepaliveTest(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create monitor and reporter
	mon := monitor.NewMonitor(cfg)
	reporter := reporter.NewReporter()

	// Start the test
	resultCh := make(chan monitor.Result, cfg.Concurrency*10)

	// Start multiple concurrent connections
	for i := 0; i < cfg.Concurrency; i++ {
		go func(id int) {
			mon.StartConnection(ctx, id, resultCh)
		}(i)
	}

	// Report results
	go func() {
		for result := range resultCh {
			reporter.Report(result)
		}
	}()

	// Wait for timeout or user interrupt
	select {
	case <-time.After(cfg.Duration):
		log.Printf("Test completed, duration: %v", cfg.Duration)
		cancel()
	case <-sigChan:
		log.Println("Received interrupt signal, stopping...")
		cancel()
	}

	// Print summary
	reporter.PrintSummary(cfg)

	return nil
}
