package main

import (
	"github.com/briandowns/spinner"
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)  // Build our new spinner
	s.Start()                                                    // Start the spinner
	time.Sleep(4 * time.Second)                                  // Run for some time to simulate work
	s.Stop()
}
