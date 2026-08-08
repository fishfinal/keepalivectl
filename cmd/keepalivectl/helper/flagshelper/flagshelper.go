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

package flagshelper

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Option defines a functional option for configuring FlagsHelper.
type Option func(*FlagsHelper)

// WithSpacing sets the spacing width for flag name alignment.
func WithSpacing(spacing int) Option {
	return func(f *FlagsHelper) {
		f.spacing = spacing
	}
}

// WithShortFlag enables or disables short flag display.
func WithShortFlag(show bool) Option {
	return func(f *FlagsHelper) {
		f.showShort = show
	}
}

// FlagsHelper provides utilities for grouping and printing command flags.
type FlagsHelper struct {
	spacing   int
	showShort bool
}

// NewFlagsHelper creates a new FlagsHelper with the given options.
func NewFlagsHelper(options ...Option) *FlagsHelper {
	f := &FlagsHelper{
		spacing:   20,
		showShort: true, // default: show short flags
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

// PrintFlagGroup prints all flags belonging to a specific group.
// If showShort is enabled, displays both short and long flag names.
func (f *FlagsHelper) PrintFlagGroup(cmd *cobra.Command, groupName string) {
	fmt.Fprintf(cmd.OutOrStdout(), "\n%s:\n", groupName)

	var flags []*pflag.Flag
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if groups, ok := flag.Annotations["group"]; ok {
			for _, g := range groups {
				if g == groupName {
					flags = append(flags, flag)
					break
				}
			}
		}
	})

	if len(flags) == 0 {
		return
	}

	for _, flag := range flags {
		f.printFlag(cmd, flag)
	}
}

// printFlag formats and prints a single flag.
func (f *FlagsHelper) printFlag(cmd *cobra.Command, flag *pflag.Flag) {
	// Build flag name with short option
	flagName := f.buildFlagName(flag)

	// Build default value string
	defaultStr := ""
	if flag.DefValue != "" && flag.DefValue != "[]" && flag.DefValue != "false" {
		defaultStr = fmt.Sprintf(" (default: %v)", flag.DefValue)
	}

	// Calculate spacing alignment (minimum 2 spaces)
	spacing := f.spacing - len(flagName)
	if spacing < 2 {
		spacing = 2
	}

	// Format and output
	fmt.Fprintf(cmd.OutOrStdout(), "  %-*s%s%s\n",
		len(flagName)+spacing,
		flagName,
		flag.Usage,
		defaultStr)
}

// buildFlagName constructs the flag name string with short option if available.
func (f *FlagsHelper) buildFlagName(flag *pflag.Flag) string {
	if !f.showShort || flag.Shorthand == "" || flag.Shorthand == "h" {
		return "--" + flag.Name
	}
	return fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
}

// PrintFlagGroupWithSeparator prints flags with a custom separator prefix.
func (f *FlagsHelper) PrintFlagGroupWithSeparator(cmd *cobra.Command, groupName string, prefix string) {
	fmt.Fprintf(cmd.OutOrStdout(), "\n%s%s:\n", prefix, groupName)

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if groups, ok := flag.Annotations["group"]; ok {
			for _, g := range groups {
				if g == groupName {
					f.printFlagWithPrefix(cmd, flag, prefix)
					break
				}
			}
		}
	})
}

// printFlagWithPrefix prints a flag with a custom prefix.
func (f *FlagsHelper) printFlagWithPrefix(cmd *cobra.Command, flag *pflag.Flag, prefix string) {
	flagName := f.buildFlagName(flag)

	defaultStr := ""
	if flag.DefValue != "" && flag.DefValue != "[]" && flag.DefValue != "false" {
		defaultStr = fmt.Sprintf(" (default: %v)", flag.DefValue)
	}

	spacing := f.spacing - len(flagName)
	if spacing < 2 {
		spacing = 2
	}

	fmt.Fprintf(cmd.OutOrStdout(), "  %s%-*s%s%s\n",
		prefix,
		len(flagName)+spacing,
		flagName,
		flag.Usage,
		defaultStr)
}

// PrintAllFlags prints all flags with their groups.
func (f *FlagsHelper) PrintAllFlags(cmd *cobra.Command, groups ...string) {
	for _, group := range groups {
		f.PrintFlagGroup(cmd, group)
	}
}

// GetFlagGroupNames returns all unique group names from flags.
func (f *FlagsHelper) GetFlagGroupNames(cmd *cobra.Command) []string {
	groupMap := make(map[string]bool)
	var groups []string

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if annotations, ok := flag.Annotations["group"]; ok {
			for _, g := range annotations {
				if !groupMap[g] {
					groupMap[g] = true
					groups = append(groups, g)
				}
			}
		}
	})

	return groups
}

// ColorizeFlagGroup prints flags with color support (if terminal supports it).
func (f *FlagsHelper) ColorizeFlagGroup(cmd *cobra.Command, groupName string, colorCode string) {
	// This is a placeholder for color support
	// You can integrate with github.com/fatih/color or similar packages
	f.PrintFlagGroup(cmd, groupName)
}
