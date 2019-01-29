package configuration

import (
	"fmt"
	"testing"
)

func TestIsValidStreamMapFormat(t *testing.T) {
	// <Stream Name>@<Cluster Name>:<Logical Server Name>
	var testSet = []struct {
		input          string
		expectedResult bool
	}{
		{"VAS_V2_SR_V2@RYD2:Server", true},
		{"VAS SR @RYD2:Server", false}, // Whitespaces not allowed between characters of the stream name
		{"1234@RYD2:Server 5", false},  // No whitespaces allowed in logical server name
		{"1234@RYD2: Server5", true},   // Whitespaces allowed after the colon
		{"1234@RYD2:Server5 ", true},   // Whitespaces allowed at the end of the mapping string
		{" 1234@RYD2:Server5", true},   // Whitespaces allowed at the beginning of the mapping string
		{"1234$RYD2:Server5", false},   // Only @ character is allowed as separator between stream name and logical server name
	}

	for _, test := range testSet {
		if result := isValidStreamMapFormat(test.input); result != test.expectedResult {
			t.Errorf("isValidStreamMapFormat(\"%s\") = %v, should be %v", test.input, result, test.expectedResult)
		}
	}
}

func BenchmarkIsValidStreamMapFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidStreamMapFormat("VAS_V2_SR_V2@RYD2:Server")
	}
}

func ExampleisValidStreamMapFormat() {
	var streamMap = "VAS_CONFIG@RYD1:Server"

	if isValidStreamMapFormat(streamMap) {
		fmt.Printf("%s is a valid format.\n", streamMap)
	} else {
		fmt.Printf("Invalid stream map format\n")
	}
}
