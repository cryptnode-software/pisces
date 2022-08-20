//go:build mage
// +build mage

package main

import (
	"github.com/cryptnode-software/pisces/lib/utility"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	utility.NewEnv(utility.NewLogger())
	return nil
}
