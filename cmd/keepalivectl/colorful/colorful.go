// colorful.go
package colorful

// ANSI standard color codes can be referenced from this project: https://github.com/fidian/ansi

import "github.com/logrusorgru/aurora"

var enableColors = true // lowercase, not directly exportable
var Colorful aurora.Aurora = aurora.NewAurora(enableColors)

// SetEnableColors allows external control of the color switch
func SetEnableColors(enable bool) {
	enableColors = enable
	Colorful = aurora.NewAurora(enableColors) // reinitialize
}

// ReviseColorfully dynamically clears color values.
func ReviseColorfully(value aurora.Value, colorfully bool) aurora.Value {
	if colorfully {
		return value
	}
	return Colorful.Reset(value)
}

// Purple is a purple color renderer.
func Purple(text string) aurora.Value {
	// ANSI standard color code, purple is 128
	return Colorful.Index(128, text)
}

// Orange is an orange color renderer.
func Orange(text string) aurora.Value {
	// ANSI standard color code, orange is 208
	return Colorful.Index(208, text)
}

// Gray is a gray color renderer.
func Gray(text string) aurora.Value {
	return Colorful.Index(7, text)
}

type SeverityColorfully func(text string) aurora.Value

// SeverityCritical is the color renderer for critical severity vulnerabilities.
var SeverityCritical SeverityColorfully = func(text string) aurora.Value {
	return Purple(text)
}

// SeverityHigh is the color renderer for high severity vulnerabilities.
var SeverityHigh SeverityColorfully = func(text string) aurora.Value {
	return Colorful.Red(text)
}

// SeverityMedium is the color renderer for medium severity vulnerabilities.
var SeverityMedium SeverityColorfully = func(text string) aurora.Value {
	return Orange(text)
}

// SeverityLow is the color renderer for low severity vulnerabilities.
var SeverityLow SeverityColorfully = func(text string) aurora.Value {
	return Colorful.Yellow(text)
}

// SeverityInfo is the color renderer for informational vulnerabilities.
var SeverityInfo SeverityColorfully = func(text string) aurora.Value {
	return Colorful.Cyan(text)
}

// SeverityUnknown is the color renderer for unknown severity vulnerabilities.
var SeverityUnknown SeverityColorfully = func(text string) aurora.Value {
	return Gray(text)
}

var _ SeverityColorfully = SeverityCritical
var _ SeverityColorfully = SeverityHigh
var _ SeverityColorfully = SeverityMedium
var _ SeverityColorfully = SeverityUnknown
var _ SeverityColorfully = SeverityLow
